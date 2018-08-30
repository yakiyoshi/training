# マイクロサービスアーキテクチャ入門（ショッピングカートサービス）

## 概要

- Go言語 1.9.x
- gin
- gorm
- Mysql

## Go環境構築

Download dependencies:

```
docker-compose run cart-app go-wrapper download
```

Start server:

```
docker-compose up -d
```

## 環境確認

http://localhost:3002

## 使い方

カート追加（1件）

```
curl -X POST \
-H "Authorization: Bearer 【JWTトークン】" \
-d 'product_id=1' \
http://localhost:3002/cart
```

カート追加（1件）for Windows PS

```
curl.exe -X POST `
-H "Authorization: Bearer 【JWTトークン】" `
-d "product_id=1" `
http://localhost:3002/cart
```

全商品取得

```
curl -X GET \
-H "Authorization: Bearer 【JWTトークン】" \
http://localhost:3002/cart
```

全商品取得 for Windows PS

```
curl.exe -X GET `
-H "Authorization: Bearer 【JWTトークン】" `
http://localhost:3002/cart
```

カート削除

```
curl -X DELETE \
-H "Authorization: Bearer 【JWTトークン】" \
-d 'cart_id=1' \
http://localhost:3002/cart
```

カート削除 for Windows PS

```
curl.exe -X DELETE `
-H "Authorization: Bearer 【JWTトークン】" `
-d "cart_id=1" `
http://localhost:3002/cart
```

カート全削除

```
curl -X DELETE \
-H "Authorization: Bearer 【JWTトークン】" \
http://localhost:3002/carts
```

カート全削除 for Windows PS

```
curl.exe -X DELETE `
-H "Authorization: Bearer 【JWTトークン】" `
http://localhost:3002/carts
```

## コンソール出力

```
docker-compose logs -f app
```
