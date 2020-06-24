import Vue from 'https://cdn.jsdelivr.net/npm/vue@2.6.11/dist/vue.esm.browser.js'

window.onload= function(){
    var search = window.location.search;
    var search_key = getSearchString('key',search);
    
    axios.get('http://167.179.81.168/bangumiSearch/'+search_key)
    .then(function(response){
        var data = response.data;
        
        for(var i = 0;i < data.length;i++){
            
            var new_str = '';
            for(var j = 0;j < data[i].staff_list.length;j++){
                if(data[i].staff_list[j] == '\''){
                    if(data[i].staff_list[j-1]==','||data[i].staff_list[j+1]==','||data[i].staff_list[j-1]==':'||data[i].staff_list[j+1]==':'||data[i].staff_list[j-1]=='{'||data[i].staff_list[j+1]=='}'){
                        new_str+='\"';continue;
                    }
                    else {new_str+='‘';continue;}
                }
                if(data[i].staff_list[j] == '\"'){new_str+='“';continue;}
                new_str+=data[i].staff_list[j];
            }

            var staff_list = JSON.parse(new_str);
            console.log(staff_list);

            var new_anime = {
                'id':data[i].bangumi_id,
                'url':data[i].cover_url,
                'score':data[i].bangumi_score,
                'desc':data[i].desc,
                'staff_list':staff_list,
                'name':data[i].name
            }
            search_result.results.push(new_anime);
        }
    })





}




var search_result = new Vue({
    el:'.page',
    data:{
        results:[],
      
    },
    methods:{
      
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

