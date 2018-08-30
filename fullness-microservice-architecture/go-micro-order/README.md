# マイクロサービスアーキテクチャ入門（商品サービス）

## 概要

- Go言語 1.9.x
- gin
- gorm
- Mysql

## Go環境構築

Download dependencies:

```
docker-compose run order-app go-wrapper download
```

Start server:

```
docker-compose up -d
```

## 環境確認

http://localhost:3001

## 使い方

注文追加

```
curl -X POST -H 'Content-Type:application/json' \
-H "Authorization: Bearer xxxxxx" \
-d '{"order_details":[{"product_id" : 1,"product_price" : 1000},{"product_id" : 2,"product_price" : 1500}]}' \
http://localhost:3003/order
```

For Windows PS
```
curl.exe -X POST -H 'Content-Type:application/json' `
-H "Authorization: Bearer <JWTトークン>" `
-d "{""""order_details"""":[{""""product_id"""" : 1,""""product_price"""" : 1000},{""""product_id"""" : 2,""""product_price"""" : 1500}]}" `
http://localhost:3003/order
```

注文リスト取得

```
curl -H 'Content-Type:application/json' \
-H "Authorization: Bearer xxxxxx" \
http://localhost:3003/orders
```

注文リスト取得 for Windows PS
```
curl.exe -X GET -H 'Content-Type:application/json' `
-H "Authorization: Bearer xxxxxx" `
http://localhost:3003/orders
```


## コンソール出力

```
docker-compose logs -f app
```

## gulpを使った自動環境再構築

下記を実行後はgoファイル or htmlファイルを修正毎に上記の再構築が実行される。

```
npm install && gulp
```