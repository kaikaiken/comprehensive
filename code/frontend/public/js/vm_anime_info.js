import Vue from 'https://cdn.jsdelivr.net/npm/vue@2.6.11/dist/vue.esm.browser.js'

window.onload = function () {

    var search = window.location.search;
    var anime_id = getSearchString('id', search);




    axios.get('http://167.179.81.168/bangumi/' + anime_id)
        .then(function (response) {

            var data = response.data[0];

            anime_info.name = data.name;
            anime_info.bangumi_id = data.bangumi_id;
            anime_info.cover_url = data.cover_url;
            anime_info.bangumi_score = data.bangumi_score;
            anime_info.vote_num = data.vote_num;
            anime_info.episode_num = data.episode_num;
            anime_info.desc = data.desc;
            anime_info.staff_list = data.staff_liste;
            anime_info.cv_list = data.cv_list;

        });




    axios.get('http://167.179.81.168:8010/CB/' + anime_id)
        .then(function (response) {
            var data = response.data;
            console.log(data);
            for (var i = 0; i < data.length; i++) {

                var new_anime = {
                    'id': data[i].bangumi_id,
                    'url': data[i].cover_url,
                    'score': data[i].bangumi_score,
                    'name': data[i].name

                }
                rec_info.animes.push(new_anime);
            }

        });


}


var anime_info = new Vue({
    el: '#anime_info',
    data: {
        name: '',
        bangumi_id: '',
        cover_url: '',
        bangumi_score: '',
        vote_num: '',
        episode_num: '',
        desc: '',
        staff_list: '',
        cv_list: ''
    },
    methods: {
        handle_fav: function () {
            if (!sessionStorage.getItem('login_status')) {
                alert("请先登录");
                window.location.href = '/index.html';
            } else {
                var username = sessionStorage.getItem('username');
                var req_id = this.bangumi_id;

                axios.get('http://47.105.129.153:1323/api/user/username/' + username + '/addBangumi/' + req_id).then(
                    function () {

                        var new_list = sessionStorage.getItem("fav_list") + ',' + req_id;
                        console.log(new_list);
                        sessionStorage.setItem("fav_list", new_list);
                        alert("success");
                    }
                );
            }
        }


    }
});



var rec_info = new Vue({
    el: '#rec_info',
    data: {
        animes: [

        ],

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