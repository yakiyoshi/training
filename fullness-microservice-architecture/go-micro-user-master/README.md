# マイクロサービスアーキテクチャ入門（顧客サービス）

## 概要

- Go言語 1.9.x
- gin
- gorm
- Mysql

## Go環境構築

Download dependencies:

```
docker-compose run user-app go-wrapper download
```

Start server:

```
docker-compose up -d
```

## 環境確認

http://localhost:3000

## 使い方

注文追加

```
curl http://localhost:3000
```

For Windows PS
```
curl.exe http://localhost:3000
```

ユーザ登録

```
curl -X POST \
-d "email=admin@example.com" \
-d "password=adminpass" \
http://localhost:3000/user
```

ユーザ登録 for Windows PS
```
curl.exe -X POST `
-d "email=admin@example.com" `
-d "password=adminpass" `
http://localhost:3000/user
```

ログイン

```
curl -X POST \
-d "email=admin@example.com" \
-d "password=adminpass" \
http://localhost:3000/login
```

ログイン for Windows PS

```
curl.exe -X POST `
-d "email=admin@example.com" `
-d "password=adminpass" `
http://localhost:3000/login
```

ユーザ情報取得

```
curl -X GET \
-H "Authorization: Bearer 【JWTトークン】" \
http://localhost:3000/me
```

ユーザ情報取得 for Windows PS

```
curl.exe -X GET `
-H "Authorization: Bearer 【JWTトークン】" `
http://localhost:3000/me
```

## コンソール出力

```
docker-compose logs -f app
```
