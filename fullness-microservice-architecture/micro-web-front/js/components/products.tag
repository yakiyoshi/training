<products>
    <!-- UIコンポーネントレイアウト -->
    <div>
        <div class="panel panel-primary">
            <div class="panel-heading">
                商品一覧
            </div>
            <div class="panel-body">
                <div class="row">
                    <div class="col-md-4 col-sm-4 col-xs-12" each={ products }>
                        <img class="img-responsive img-portfolio img-hover" src="{ product_image }" alt="">
                        <h4>商品名:{ product_name } 価格:{ product_price }円</h4>
                        <button if={ is_login } onclick={ parent.cart_in } class="btn btn-primary btn-block">カートに入れる</button>
                    </div>
                </div>
            </div>
        </div>
    </div>
    
    <!-- UIコンポーネントロジック -->
    <script>
    //-------------------------
    // プロパティ
    //-------------------------
    this.products = []
    this.is_login = false
    var self = this
    //-------------------------
    // 商品一覧取得リクエスト
    //-------------------------
    this.fetch_carts = function(){
        const method = "GET"
        const headers = {'Accept': 'application/json'}
        fetch("http://localhost:3001/products", {method, headers})
            .then(function (resp) { return resp.json() })
            .then(function (json) {
                if (json) {
                    console.log(json)
                    self.products = json
                    self.update()
                    observer.trigger("products_update", self.products)
                }
            })
    }
    self.fetch_carts()
    //-------------------------
    // 商品をカートに入れる
    //-------------------------
    this.cart_in = function(event) {
        var product = event.item
        var token = Cookies.get('token')
        if(token){
            const method = "POST"
            const headers = {
                'Accept': 'application/json',
                'Content-Type': 'application/x-www-form-urlencoded; charset=utf-8',
                'Authorization': 'Bearer ' + token
            }
            const obj = {product_id:product.ID}
            const body = Object.keys(obj).map((key)=>key+"="+encodeURIComponent(obj[key])).join("&")
            fetch("http://localhost:3002/cart", {method, headers,body})
                .then(function (resp) { return resp.json() })
                .then(function (json) {
                    if (json) {
                        console.log(json)
                        observer.trigger("cart_in", "");
                    }
                })
        }
    }
    //-------------------------
    // ログインイベント受信時
    //-------------------------
    observer.on("login", function(user_id) {
        self.is_login = true
        self.update()
    })
    //-------------------------
    // ログアウトイベント受信時
    //-------------------------
    observer.on("logout", function(user_id) {
        self.is_login = false
        self.update()
    })
    </script>

    <!-- UIコンポーネントレイアウトデザイン -->
    <style>
        :scope
        em{ font-size: 1rem; color:#f00; }
    </style>
</products>
