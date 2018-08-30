<login>
    <!-- UIコンポーネントレイアウト -->
    <div>
        <form if={ is_login == false } onsubmit={ login }>
            <input ref="username" />
            <input ref="password" />
            <button ref="submit" class="btn btn-success">ログイン</button>
        </form>
        <div if={ is_login }>
            <strong>ようこそ！{ email }さん</strong>
            <form onsubmit={ logout }>
                <button ref="logout" class="btn btn-success">ログアウト</button>
            </form>
        </div>
    </div>
    <!-- UIコンポーネントロジック -->
    <script>
    //-------------------------
    // プロパティ
    //-------------------------
    this.is_login = false
    this.user_id = ''
    var self = this
    //-------------------------
    // トークンチェック
    //-------------------------
    var token = Cookies.get('token')
    if(token){
        fetch_user_info(token)
    }
    //-------------------------
    // ログインリクエスト
    //-------------------------
    this.login = function(e){
        e.preventDefault()
        var username = this.refs.username.value
        var password = this.refs.password.value
        const obj = {user_id: username,password:password}
        const method = "POST"
        const body = Object.keys(obj).map((key)=>key+"="+encodeURIComponent(obj[key])).join("&")
        const headers = {'Accept': 'application/json', 'Content-Type': 'application/x-www-form-urlencoded; charset=utf-8'}
        fetch("http://localhost:3000/login", {method, headers, body})
            .then(function (resp) {
                return resp.json()
            })
            .then(function (json) {
                if (json.token) {
                    Cookies.set("token",json.token)
                    self.is_login = true
                    self.update()
                    fetch_user_info(json.token)
                }
            })
    }
    //-------------------------
    // ユーザ情報取得リクエスト
    //-------------------------
    function fetch_user_info(token){
        const method = "GET"
        const headers = {
            'Accept': 'application/json',
            'Content-Type': 'application/x-www-form-urlencoded; charset=utf-8',
            'Authorization': 'Bearer ' + token
        }
        fetch("http://localhost:3000/me", {method, headers})
            .then(function (resp) {
                return resp.json()
            })
            .then(function (json) {
                if (json.email) {
                    self.is_login = true
                    self.email = json.email
                    self.user_id = json.user_id
                    self.update()
                    // ログインしたイベントを送出
                    observer.trigger("login", self.user_id);
                }
            })
    }
    //-------------------------
    // ログアウトリクエスト
    //-------------------------
    this.logout = function(e){
        Cookies.remove('token') // クッキー削除
        self.is_login = false
        self.update()
        observer.trigger("logout", "");
    }
    </script>

    <!-- UIコンポーネントレイアウトデザイン -->
    <style>
        :scope
        form{ display: inline }
    </style>
</login>
