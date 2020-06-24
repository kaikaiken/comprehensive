import Vue from 'https://cdn.jsdelivr.net/npm/vue@2.6.11/dist/vue.esm.browser.js'


window.onload = function () {




    navbar.login_status = sessionStorage.getItem("login_status");

    var search = window.location.search;
    var page_key = getSearchString('page', search);
    index_rec.cur_page = page_key > 0 ? parseInt(page_key) : 1;


    axios.get('http://47.105.129.153:8000/bangumiAll/' + index_rec.cur_page)
        .then(function (response) {
            var data = response.data;
            console.log(data);
            for (var i = 0; i < data.length; i++) {

                var new_anime = {
                    name: data[i].name,
                    id: data[i].bangumi_id,
                    cover_url: data[i].cover_url,
                    score: data[i].bangumi_score
                }
                index_rec.animes.push(new_anime);
            }

        });


}





var index_rec = new Vue({
    el: '#index_rec',
    data: {
        cur_page: 4,
        total_page: 223,
        animes: [
        ]
    },
    methods: {
        handle_fav: function (event) {

            if (!sessionStorage.getItem('login_status')) navbar.login_show = 1;
            else {
                var username = sessionStorage.getItem('username');
                var el_id = event.target.id;
                var req_id = el_id.substring(3);
                axios.get('http://47.105.129.153:1323/api/user/username/' + username + '/addBangumi/' + req_id).then(
                    function () {
                        var new_list = sessionStorage.getItem("fav_list") + ',' + req_id;
                        sessionStorage.setItem("fav_list", new_list);
                    }
                );
            }



        }
    }
})





var navbar = new Vue({
    el: '#navbar',
    data: {
        login_status: 0,
        login_show: 0,
        user_name: '',
        pass_word: '',
        success_username: ''
    },
    methods: {
        handleClick: function () {

            if (!this.login_status || this.login_status == 0) {
                this.login_show = 1;
            } else {
                window.location.href = "/my_page.html?name=" + sessionStorage.getItem("username");
            }






        },
        login_request: function () {

            var p_this = this;

            var data = new FormData();
            data.append('username', this.user_name);
            data.append('password', this.pass_word);

            axios.post('http://47.105.129.153:1323/api/user/login', data)
                .then(
                    function (response) {

                        axios.get('http://47.105.129.153:1323/api/user/username/' + p_this.user_name).then(
                            function (response) {
                                var data = response.data;
                                sessionStorage.setItem("login_status", 1);
                                sessionStorage.setItem("username", data.username);
                                sessionStorage.setItem("avatar", data.avatar);
                                sessionStorage.setItem("email", data.email);
                                sessionStorage.setItem("fav_list", data.bangumi_list);
                                p_this.login_status = 1;
                                p_this.login_show = 0;
                                window.location.href = "my_page.html?name=" + p_this.user_name;
                            }
                        );








                    }).catch(function (error) {
                        if (error.response) console.log(error);
                    })





        }
    }
});







function getSearchString(key, Url) {
    var str = Url;
    str = str.substring(1, str.length); // 获取URL中?之后的字符（去掉第一位的问号）

    var arr = str.split("&");
    var obj = new Object();
    // 将每一个数组元素以=分隔并赋给obj对象
    for (var i = 0; i < arr.length; i++) {
        var tmp_arr = arr[i].split("=");
        obj[decodeURIComponent(tmp_arr[0])] = decodeURIComponent(tmp_arr[1]);
    }
    return obj[key];
}