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
    Admin.AddResource(&Order{})
    Admin.AddResource(&OrderDetail{})
    // HTTPリクエストマルチプレクサ作成
    mux := http.NewServeMux()
    // 管理画面をマルチプレクサにマウント
    Admin.MountTo("/admin", mux)
    fmt.Println("注文サービスDB管理画面起動 PORT:9003")
    go http.ListenAndServe(":9003", mux) // 並行処理で起動
}

//------------------------------
// Database
//------------------------------
var db *gorm.DB // グローバル変数としてDBオブジェクトを保持
// 注文エンティティ
type Order struct {
    gorm.Model
    UserId int                 `json:"user_id"`
    OrderDetails []OrderDetail `gorm:"foreignkey:OrderId";json:"order_details"`
}
// 注文明細エンティティ
type OrderDetail struct {
    gorm.Model
    OrderId int      `json:"order_id"`
    ProductId int    `json:"product_id"`
    ProductPrice int `json:"product_price"`
}
// DB接続
func gormConnect() *gorm.DB {
    DBMS     := "mysql"
    USER     := "root"
    PASS     := "mysql"
    PROTOCOL := "tcp(order-mysql:3306)"
    DBNAME   := "micro_order"
    // DB接続
    CONNECT := USER+":"+PASS+"@"+PROTOCOL+"/"+DBNAME+"?charset=utf8&parseTime=true"
    db,err := gorm.Open(DBMS, CONNECT)
    if err != nil {
        panic(err.Error())
    } else {
        fmt.Println("DB接続成功")
    }
    // テーブル作成
    if !db.HasTable(&Order{}) {
        db.CreateTable(&Order{})
    }
    if !db.HasTable(&OrderDetail{}) {
        db.CreateTable(&OrderDetail{})
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
    config.AllowHeaders = []string{"Authorization"}
    r.Use(cors.New(config))
    // ルーティング設定
    r.GET("/", IndexHandler)
    r.GET("/orders", GetAllOrderItemHandler)
    r.POST("/order", CreateOrderItemHandler)

    r.Run(":3003")
    fmt.Println("注文サービス起動完了")
}

// サービストップ
func IndexHandler(c *gin.Context) {
    c.JSON(200, gin.H{"message": "注文サービスへようこそ!!"})
}
// 注文全取得
func GetAllOrderItemHandler(c *gin.Context) {
    // トークン検証
    user,err := validateJwtToken(c.Request)
    if err != nil{
        c.JSON(400,ErrorResponse{1,"ユーザ認証していません"})
    }
    orders := []Order{}
    db.Where(map[string]interface{}{"user_id": user.UserId}).Find(&orders)
    for i:= 0; i<len(orders); i++{
        db.Model(&orders[i]).Association("OrderDetails").Find(&orders[i].OrderDetails)
    }
    if len(orders) > 0{
        c.JSON(200,orders)
    } else {
        c.JSON(400,ErrorResponse{1,"アイテムが存在しません"})
    }
}
// 注文追加
func CreateOrderItemHandler(c *gin.Context) {
    // トークン検証
    user,err1 := validateJwtToken(c.Request)
    if err1 != nil{
        c.JSON(400,ErrorResponse{1,"ユーザ認証していません"})
    }
    var orderJson OrderDetailArrayRequest
    err2 := c.BindJSON(&orderJson)
    if err2 != nil {
        c.JSON(400,ErrorResponse{1,"注文明細がパースできませんでした"})
        return
    }
    fmt.Printf("注文明細パース結果: %#v\n", orderJson.OrderDetails)
    order := Order{UserId:user.UserId, OrderDetails:orderJson.OrderDetails}
    db.NewRecord(order)
    db.Create(&order)
    db.Save(&order)
	 // メッセージ送信
	 jsonbyte, err3 := json.Marshal(user)
	 if err3 == nil{
		 sentMessage("order-complete", jsonbyte)
	 }
    c.JSON(200,order)
}
// 注文明細レスポンス用型
type OrderDetailArrayRequest struct {
    OrderDetails []OrderDetail `json:"order_details" binding:"required"`
}
// エラーレスポンス用型
type ErrorResponse struct{
    ErrorCode int  `json:"error_code"`
    Message string `json:"error_message"`
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
// RabbitMQ Sending Setting
//------------------------------
// RabbitMQ接続情報
var amqpURI string = "amqp://user:bitnami@rabbitmq:5672"
// エラーハンドリング
func failOnError(err error, msg string) {
    if err != nil {
        panic(fmt.Sprintf("%s: %s", msg, err))
    }
}
// メッセージ送信
func sentMessage(queueName string, payload []byte){
    conn, err := amqp.Dial(amqpURI)
    failOnError(err, "RabbitMQ接続失敗")
    defer conn.Close()
    channel, err1 := conn.Channel()
    failOnError(err1, "Failed to open a channel")
    err2 := channel.Publish(
        "",         // exchange
        queueName,  // routing key
        false,      // mandatory
        false,      // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body: payload,
        })
    failOnError(err2, "メッセージ送信失敗")
}
