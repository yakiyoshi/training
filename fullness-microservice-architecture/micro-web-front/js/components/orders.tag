<orders>
    <!-- UIコンポーネントレイアウト -->
    <div>
        <div class="panel panel-warning">
            <div class="panel-heading">
                注文一覧
            </div>
            <div class="panel-body">
                <table if={ orders != 0 } class="table">
                    <tr>
                        <th>注文日時</th>
                        <th>商品数</th>
                        <th>金額</th>
                    </tr>
                    <tr each={ orders }>
                        <td>{ order_date }</td>
                        <td>{ order_item_num }</td>
                        <td>{ order_total_cost }</td>
                    </tr>
                </table>
                <div if={ orders == 0 }>
                    注文はありません
                </div>
            </div>
        </div>
    </div>
    
    <!-- UIコンポーネントロジック -->
    <script>
    //-------------------------
    // プロパティ
    //-------------------------
    this.orders = []
    var token = Cookies.get('token')
    var self = this
    //-------------------------
    // 注文リスト取得
    //-------------------------
    this.fetch_orders = function() {
        const method = "GET"
        const headers = {
            'Accept': 'application/json',
            'Content-Type': 'application/x-www-form-urlencoded; charset=utf-8',
            'Authorization': 'Bearer ' + token
        }
        fetch("http://localhost:3003/orders", {method, headers})
            .then(function (resp) { return resp.json() })
            .then(function (json) {
                if (json) {
                    console.log(json)
                    self.orders = []
                    json.forEach( function( order ) {
                        var order_total_cost = 0
                        order.OrderDetails.forEach( function( order_detail ) {
                            // 合計金額計算
                            order_total_cost += order_detail.product_price
                        })
                        self.orders.push({
                            order_date : order.CreatedAt,
                            order_item_num : order.OrderDetails.length,
                            order_total_cost : order_total_cost
                        })
                    })
                    self.update()
                }
            })
    }
    // クッキーある場合は初期処理でカートアイテム取得
    if(token){ self.fetch_orders() }
    //-------------------------
    // イベント処理
    //-------------------------
    // 注文確定時
    observer.on("order_complete", function(e) {
        self.fetch_orders()
    })
    // ログイン時
    observer.on("login", function(user_id) {
        token = Cookies.get('token')
        self.fetch_orders()
    })
    // ログアウト時
    observer.on("logout", function(user_id) {
        self.orders = []
        self.update()
    })
    </script>

    <!-- UIコンポーネントレイアウトデザイン -->
    <style>
        :scope
        em{ font-size: 1rem; color:#f00; }
    </style>
</orders>