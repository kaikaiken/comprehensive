import Vue from 'https://cdn.jsdelivr.net/npm/vue@2.6.11/dist/vue.esm.browser.js'

window.onload = function(){
    if(sessionStorage.getItem("login_status")==0)window.location.href = "/index.html";

    var search = window.location.search;

    var page_key = getSearchString("page",search);
    my_fav.cur_page=page_key>0?parseInt(page_key):1;


    my_fav.user_name = sessionStorage.getItem("username");
    my_fav.email = sessionStorage.getItem("email");
    my_fav.email = sessionStorage.getItem("avatar");


    var fav_list = sessionStorage.getItem("fav_list");
    var fav_array = fav_list.split(',');

    axios.get('http://167.179.81.168/CF/'+fav_list).then(function(response){
        var data = response.data;
        for(let i = 0;i < data.length;i++){
            var new_anime={
                'id':data[i].bangumi_id,
                'url':data[i].cover_url,
                'score':data[i].bangumi_score,
                'name':data[i].name
            }
            my_fav.rec_list.push(new_anime);
        }
    });

    var request_start = (my_fav.cur_page - 1) * 12;
    var request_end = Math.min(fav_array.length,my_fav.cur_page * 12);


    for(var i = 0;i < fav_array.length;i++){
        
axios.get('http://167.179.81.168/bangumi/'+fav_array[i])
.then(function(response){
 
    var data = response.data[0];
    
    var new_anime = {
        'name': data.name,
        'score':data.bangumi_score,
        'url':data.cover_url,
        'id':data.bangumi_id,
        'ep_num':data.episode_num
    }
    my_fav.fav_list.push(new_anime);

});











    }





}




var my_fav = new Vue({
    el:'#my_fav',
    data:{
        user_name:'',
        fav_list:[],
        rec_list:[],
        email:'',
        avatar:'',
        total_page:1,
        cur_page:1
    },
    methods:{
        handle_del:function(target){
                var username = sessionStorage.getItem('username');
                var el_id = event.target.id;
                var req_id = el_id.substring(3);
                console.log(req_id);
              
        },
        handle_logout:function(){
            sessionStorage.clear();
            sessionStorage.setItem('login_status',0);
            window.location.href="/index.html";
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