<cart>
    <!-- UIコンポーネントレイアウト -->
    <div>
        <div class="panel panel-success">
            <div class="panel-heading">
                買い物カゴ
            </div>
            <div class="panel-body">
                <table if={ cart_items != 0 } class="table">
                    <tr>
                        <th>商品名</th>
                        <th>価格</th>
                    </tr>
                    <tr each={ cart_items }>
                        <td>{ product_name }</td>
                        <td>{ product_price }円</td>
                    </tr>
                    <tr>
                        <td colspan="2">
                            <div class="text-center"><b>合計金額: { total_cost }円</b></div>
                        </td>
                    </tr>
                    <tr>
                        <td colspan="2">
                            <button onclick={ buy } class="btn btn-warning btn-block">購入</button>
                        </td>
                    </tr>
                </table>
                <div if={ cart_items == 0 }>
                    買い物カゴに商品はありません
                </div>
            </div>
        </div>
    </div>
    
    <!-- UIコンポーネントロジック -->
    <script>
    //-------------------------
    // プロパティ
    //-------------------------
    this.cart_items = []
    this.products = []
    this.total_cost = 0
    var token = Cookies.get('token')
    var self = this
    //-------------------------
    // カートアイテム取得
    //-------------------------
    this.fetch_cart_items = function() {
        const method = "GET"
        const headers = {
            'Accept': 'application/json',
            'Content-Type': 'application/x-www-form-urlencoded; charset=utf-8',
            'Authorization': 'Bearer ' + token
        }
        fetch("http://localhost:3002/cart", {method, headers})
            .then(function (resp) { return resp.json() })
            .then(function (json) {
                if (json) {
                    self.cart_items = []
                    self.total_cost = 0
                    console.log(json)
                    json.forEach( function( cart_item ) {
                        self.products.forEach( function( product ) {
                            if(product.ID == cart_item.product_id){
                                self.cart_items.push({
                                    product_name : product.product_name,
                                    product_price : product.product_price,
                                })
                            }
                        })
                    })
                    // 合計金額
                    self.cart_items.forEach(function(item){
                        self.total_cost += item.product_price
                    })
                    self.update()
                }
            })
    }
    //-------------------------
    // 商品購入
    //-------------------------
    this.buy = function () {
        const method = "POST"
        const headers = {
            'Accept': 'application/json',
            'Content-Type': 'application/x-www-form-urlencoded; charset=utf-8',
            'Authorization': 'Bearer ' + token
        }
        const body = JSON.stringify({ order_details : self.cart_items })
        fetch("http://localhost:3003/order", {method, headers, body})
            .then(function (resp) { return resp.json() })
            .then(function (json) {
                if (json) {
                    self.delete_cart_items() //カート削除
                    observer.trigger("order_complete", "");
                }
            })
    }
    //-------------------------
    // カートアイテム全削除
    //-------------------------
    this.delete_cart_items = function () {
        const method = "DELETE"
        const headers = {
            'Accept': 'application/json',
            'Content-Type': 'application/x-www-form-urlencoded; charset=utf-8',
            'Authorization': 'Bearer ' + token
        }
        fetch("http://localhost:3002/carts", {method, headers })
                .then(function (resp) { return resp.json() })
                .then(function (json) {
                    if (json) {
                        self.fetch_cart_items() //カート更新
                    }
                })
    }
    // クッキーある場合は初期処理でカートアイテム取得
    if(token){
        self.fetch_cart_items()
    }
    //-------------------------
    // イベント処理
    //-------------------------
    // 商品情報更新時
    observer.on("products_update", function(products) {
        self.products = products
    })
    // カートにいれた時
    observer.on("cart_in", function(e) {
        self.fetch_cart_items()
    })
    // ログイン時
    observer.on("login", function(user_id) {
        token = Cookies.get('token')
        self.fetch_cart_items()
    })
    // ログアウト時
    observer.on("logout", function(user_id) {
        self.cart_items = []
        self.update()
    })
    </script>

    <!-- UIコンポーネントレイアウトデザイン -->
    <style>
        :scope
        em{ font-size: 1rem; color:#f00; }
    </style>
</cart>