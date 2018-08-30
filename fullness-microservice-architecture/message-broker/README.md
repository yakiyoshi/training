# マイクロサービスアーキテクチャ入門（メッセージキューサービス）

注文追加

```
curl -X POST -H 'Content-Type:application/json' \
-H "Authorization: Bearer JWTトークン" \
-d '{"order_details":[{"product_id" : 1,"product_price" : 1000},{"product_id" : 2,"product_price" : 1500}]}' \
http://localhost:3003/order
```

注文追加 for Windows PS

```
curl.exe -X POST -H 'Content-Type:application/json' `
-H "Authorization: Bearer JWTトークン" `
-d "{""""order_details"""":[{""""product_id"""" : 1,""""product_price"""" : 1000},{""""product_id"""" : 2,""""product_price"""" : 1500}]}" `
http://localhost:3003/order
```