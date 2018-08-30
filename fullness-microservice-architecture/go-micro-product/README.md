# マイクロサービスアーキテクチャ入門（製品サービス）

## 概要

- Go言語 1.9.x
- gin
- gorm
- Mysql

## Go環境構築

Download dependencies:

```
docker-compose run product-app go-wrapper download
```

Start server:

```
docker-compose up -d
```

## 環境確認

商品カタログサービストップ

```
curl http://localhost:3001
```

商品カタログサービストップ for Windows PS

```
curl.exe http://localhost:3001
```

全商品取得
```
curl http://localhost:3001/products
```

全商品取得 for Windows PS
```
curl.exe http://localhost:3001/products
```
