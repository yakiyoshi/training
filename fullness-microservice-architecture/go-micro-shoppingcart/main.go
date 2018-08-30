package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "gopkg.in/resty.v1"
    "encoding/json"
    "net/http"
    "strconv"
	_ "github.com/qor/qor"
    "github.com/qor/admin"
    "github.com/streadway/amqp"
)

//------------------------------
// EntryPoint
//------------------------------
func main() {
	// DB接続
    db = gormConnect()
    // メッセージ受信
    go receiveMessage()
	// Admin画面
    launchAdmin()
	// RestAPI設定
    launchRestApi()
}

//------------------------------
// Admin UI
//------------------------------
func launchAdmin(){
    // 管理画面初期化
    Admin := admin.New(&admin.AdminConfig{DB: db})
    // 管理対象のgormテーブルを指定
    Admin.AddResource(&Cart{})
    // HTTPリクエストマルチプレクサ作成
    mux := http.NewServeMux()
    // 管理画面をマルチプレクサにマウント
    Admin.MountTo("/admin", mux)
    fmt.Println("ショッピングサービスDB管理画面起動 PORT:9002")
    go http.ListenAndServe(":9002", mux) // 並行処理で起動
}

//------------------------------
// Database
//------------------------------
var db *gorm.DB // グローバル変数としてDBオブジェクトを保持
// エンティティ
type Cart struct {
    gorm.Model
    UserId        int `json:"user_id"`
    ProductId     int `json:"product_id"`
}

// DB接続
func gormConnect() *gorm.DB {
    DBMS     := "mysql"
    USER     := "root"
    PASS     := "mysql"
    PROTOCOL := "tcp(cart-mysql:3306)"
    DBNAME   := "micro_cart"
    // DB接続
    CONNECT := USER+":"+PASS+"@"+PROTOCOL+"/"+DBNAME+"?charset=utf8&parseTime=true"
    db,err := gorm.Open(DBMS, CONNECT)
    if err != nil {
        panic(err.Error())
    } else {
        fmt.Println("DB接続成功")
    }
    // テーブル作成
    if !db.HasTable(&Cart{}) {
        // テーブル作成
        db.CreateTable(&Cart{})
    }
    return db
}

//------------------------------
// REST API
//------------------------------
// REST API起動
func launchRestApi(){
    r := gin.Default()
    // CORS設定
    config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    config.AllowMethods = []string{"GET", "POST", "DELETE"}
    config.AllowHeaders = []string{"Authorization"}
    r.Use(cors.New(config))
    // ルーティング設定
    r.GET("/", IndexHandler)
    r.GET("/cart", GetAllCartItemHandler)
    r.POST("/cart", CreateCartItemHandler)
    r.DELETE("/cart", DeleteCartItemHandler)
    r.DELETE("/carts", DeleteAllCartItemHandler)

    r.Run(":3002")
    fmt.Println("ショッピングカートサービス起動完了")
}

// エラーレスポンスオブジェクト定義
type ErrorResponse struct{
    ErrorCode int `json: error_code`
    Message string `json error_message`
}

// サービルルートエントリーポイント
func IndexHandler(c *gin.Context) {
    c.JSON(200, gin.H{"message": "ショッピングカタログサービスへようこそ!!"})
}

// カートにアイテム取得（ユーザ別）
func GetAllCartItemHandler(c *gin.Context) {
    // トークン検証
    user,err := validateJwtToken(c.Request)
    if err != nil{
        c.JSON(400,ErrorResponse{1,"ユーザ認証していません"})
    }
    cartContent := []Cart{}
    db.Where(map[string]interface{}{"user_id": user.UserId}).Find(&cartContent)
    if len(cartContent) > 0{
        c.JSON(200,cartContent)
    } else {
        c.JSON(200,[]Cart{})
    }
}

// カートにアイテム追加（ユーザ別）
func CreateCartItemHandler(c *gin.Context) {
    // トークン検証
    user,err := validateJwtToken(c.Request)
    if err != nil{
        c.JSON(400,ErrorResponse{1,"ユーザ認証していません"})
    }
    productId,err := strconv.Atoi(c.PostForm("product_id"))
    if err != nil{
        c.JSON(400,ErrorResponse{2,"リクエストエラー"})
    }
    cartItem := Cart{
        ProductId: productId,
        UserId:user.UserId,
    }
    db.NewRecord(cartItem)
    db.Create(&cartItem)
    db.Save(&cartItem)
    c.JSON(200,cartItem)
}

// カートにアイテム削除（ユーザ別）
func DeleteCartItemHandler(c *gin.Context){
    // トークン検証
    user,err := validateJwtToken(c.Request)
    if err != nil{
        c.JSON(400,ErrorResponse{1,"ユーザ認証していません"})
    }
    cartItemId,_ := strconv.Atoi(c.PostForm("cart_id"))
    cart := Cart{}
    cart.ID = uint(cartItemId)
    cart.UserId = user.UserId
    db.First(&cart)
    db.Delete(&cart)
    c.JSON(200, gin.H{"message": "カートアイテム削除しました"})
}

// カートアイテム全削除（ユーザ別）
func DeleteAllCartItemHandler(c *gin.Context){
    // トークン検証
    user,err := validateJwtToken(c.Request)
    if err != nil{
        c.JSON(400,ErrorResponse{1,"ユーザ認証していません"})
    }
    cartContent := []Cart{}
    db.Where(map[string]interface{}{"user_id": user.UserId}).Find(&cartContent)
    db.Delete(&cartContent)
    c.JSON(200, gin.H{"message": "ユーザID:" + user.Email + "のカートアイテムを全削除しました"})
}

//------------------------------
// 認証関連
//------------------------------
type User struct{
    UserId int   `json:"user_id"`
    Email string `json:"email"`
}

func validateJwtToken(req *http.Request) (*User,error){
    // 顧客サービスにjwt検証リクエスト
    fmt.Println("リクエストトークン確認:" + req.Header.Get("Authorization"))
    resp, err := resty.R().
        SetHeader("Authorization", req.Header.Get("Authorization")).
        Get("http://user-app:3000/me")
    if err != nil{
        fmt.Println("認証リクエストエラー:", err)
        return nil,err
    }
    // エラーレスポンス時
    if resp.StatusCode() != 200{
        fmt.Println("認証時エラーレスポンス:", err)
        return nil,err
    }
    user := new(User)
    fmt.Println("認証結果:" + resp.String())
    if err := json.Unmarshal(([]byte)(resp.String()), user); err != nil {
        fmt.Println("JSONパースエラー:", err)
        return nil,err
    }
    fmt.Println(user)
    return user,nil
}

//------------------------------
// RabbitMQ Receive Setting
//------------------------------
var amqpURI string = "amqp://user:bitnami@rabbitmq:5672"

func failOnError(err error, msg string) {
    if err != nil {
        panic(fmt.Sprintf("%s: %s", msg, err))
    }
}

func receiveMessage(){
    conn, err := amqp.Dial(amqpURI)
    failOnError(err, "RabbitMQ接続失敗")
    defer conn.Close()

    channel, err := conn.Channel()
    failOnError(err, "チャンネルオープン失敗")

    q, err := channel.QueueDeclare(
        "order-complete", // name
        false,      // durable
        false,      // delete when unused
        false,      // exclusive
        false,      // no-wait
        nil,        // arguments
    )
    failOnError(err, "キュー宣言失敗")

    messages, err := channel.Consume(
        q.Name,     // queue
        "",         // consumer
        true,       // auto-ack
        false,      // exclusive
        false,      // no-local
        false,      // no-wait
        nil,        // arguments
    )
    failOnError(err, "受信設定失敗")

    forever := make(chan bool)
    // 並行処理でメッセージ受信
    go func() {
        for data := range messages {
            fmt.Printf("%s\n", data.Body)
            user := new(User)
            if err := json.Unmarshal(data.Body, user); err != nil {
                failOnError(err, "JSONパースエラー")
            }
            // ユーザのカート削除処理
            cartContent := []Cart{}
            db.Where(map[string]interface{}{"user_id": user.UserId}).Find(&cartContent)
            db.Delete(&cartContent)
            fmt.Println("カート削除処理完了:")
        }
    }()
    fmt.Printf("メッセージ受信開始 To exit press CTRL+C\n")
    <-forever
}
