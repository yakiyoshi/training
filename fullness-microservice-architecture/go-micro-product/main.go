package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"net/http"
    _ "github.com/qor/qor"
    "github.com/qor/admin"
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
    Admin.AddResource(&Product{})
    // HTTPリクエストマルチプレクサ作成
    mux := http.NewServeMux()
    // 管理画面をマルチプレクサにマウント
    Admin.MountTo("/admin", mux)
    fmt.Println("顧客サービスDB管理画面起動 PORT:9001")
    go http.ListenAndServe(":9001", mux) // 並行処理で起動
}
//------------------------------
// Database
//------------------------------
var db *gorm.DB // グローバル変数としてDBオブジェクトを保持
// エンティティ
type Product struct {
    gorm.Model
    ProductId     int    `json:"product_id"`
    ProductName   string `json:"product_name"`
    ProductImage  string `json:"product_image"`
    ProductPrice  int    `json:"product_price"`
}

// DB接続
func gormConnect() *gorm.DB {
    DBMS     := "mysql"
    USER     := "root"
    PASS     := "mysql"
    PROTOCOL := "tcp(product-mysql:3306)"
    DBNAME   := "micro_product"
    // DB接続
    CONNECT := USER+":"+PASS+"@"+PROTOCOL+"/"+DBNAME+"?charset=utf8&parseTime=true"
    db,err := gorm.Open(DBMS, CONNECT)
    if err != nil {
        panic(err.Error())
    } else {
        fmt.Println("DB接続成功")
	}
	// テーブル作成
    if !db.HasTable(&Product{}) {
        // テーブル作成
        db.CreateTable(&Product{})
        // レコード追加
        product1 := Product{
            ProductId:1,
            ProductName: "デジカメ",
            ProductImage:"https://www.image-pit.com/img/ms/camera.png",
            ProductPrice:50000,
        }
        db.NewRecord(product1)
        db.Create(&product1)
        db.Save(&product1)
        product2 := Product{
            ProductId:2,
            ProductName: "パソコン",
            ProductImage:"https://www.image-pit.com/img/ms/pc.png",
            ProductPrice:100000,
        }
        db.NewRecord(product2)
        db.Create(&product2)
        db.Save(&product2)
        product3 := Product{
            ProductId:3,
            ProductName: "スマートフォン",
            ProductImage:"https://www.image-pit.com/img/ms/smartphone.png",
            ProductPrice:70000,
        }
        db.NewRecord(product3)
        db.Create(&product3)
        db.Save(&product3)
    }
    return db
}

//------------------------------
// REST API
//------------------------------
// REST API起動
func launchRestApi(){
    r := gin.Default()
	// CORS設定(クロス」オリジンソースシェアリング)
	// クロスオリジン
	// プロトコル(http)・ホスト(localhost)・ポート(3301)　→　オリジン
	// 通常のブラウザは上記以外のオリジンからJSでリクエストなげるとブラウザ側ではじく（クロスオリジン制約）
    config := cors.DefaultConfig()
    config.AllowAllOrigins = true // HTTPヘッダーを追加（クロスオリジン制約を解除している）
    config.AllowHeaders = []string{"Authorization"}
    r.Use(cors.New(config))
    // ルーティング設定
    r.GET("/", IndexHandler)
    r.GET("/products", GetAllProductsHandler) // 商品全取得
    r.Run(":3001")
    fmt.Println("商品カタログサービス起動完了")
}

// エラーレスポンスオブジェクト定義
type ErrorResponse struct{
    ErrorCode int `json: error_code`
    Message string `json error_message`
}

// サービルルートエントリーポイント
func IndexHandler(c *gin.Context) {
    c.JSON(200, gin.H{"message": "商品カタログサービスへようこそ!!"})
}

// 商品全取得
func GetAllProductsHandler(c *gin.Context) {
    products := []Product{}
    db.Find(&products)
    if len(products) > 0{
        c.JSON(200,products)
    } else {
        c.JSON(400,ErrorResponse{1,"アイテムが存在しません"})
    }
}