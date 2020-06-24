# 番剧推荐系统

 ## Backend, Golang

### 相关工具

- 编程语言: [Golang](https://golang.org)

- 使用框架: [go-kit](https://github.com/go-kit/kit)

- 重要额外工具/库:
  -  [gRPC-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
  - [protobuf](https://github.com/protocolbuffers/protobuf)
  - [protoc-gen-go](https://github.com/golang/protobuf)

- 模块管理: [go mod](https://github.com/golang/go/wiki/Modules)

### 使用方法

运行三个微服务

```sh
cd info
go run info.go
```

```sh
cd log
go run recommend.go
```

```
cd user
go run user.go
```

运行网关

```sh
cd gateway
go run gateway
```

### API文档

后端使用了gRPC, 所以并没有暴露在外的RESTful API, 所以使用了gRPC-gateway帮助翻译http json为gRPC req与rsp

#### user

user部分主要包括用户注册与登录以及修改用户信息

##### 注册

- URL:  post: "/api/v1/users"

- Req body:  

  ```json
  {
  	"password" : string,
  	"username" : string,
  	"phone" : string
  }
  ```

- Rsp body

  ```json
  {
  	"password" : int32,
  	"msg" : string
  }
  ```

  

##### 登录

- URL: post: "/api/v1/sessions"

- Req body

  ```json
  {
  	"id": int32,
  	"password": string
  }
  ```

- Rsp body

  ```json
  {
  	"id" : int32,
  	"password" : string,
  	"status" : bool,
  	"phone" : string,
  	"msg" : string,
  	"jwt" : string
  }
  ```

##### 个人信息

- URL:  post: "/api/v1/profiles"

- Req body

  ```json
  {
  	"id": int32,
  	"password": string,
  	"username": string,
  	"status": bool,
  	"phone": string,
  	"email": string
  }
  ```
  
- Rsp body

  ```json
  {
  	"ok" :bool,
  	"msg" :string,
  	"id": int32,
  	"password": string,
  	"username": string,
  	"status": bool,
  	"phone": string,
  	"email": string,
    "jwt": string
  }
  ```

  

#### recommend

##### 收藏

- URL:  get: "/api/v1/{id}/favorites"

- Req body

  ```josn
  nil
  ```

- Rsp body

  ```json
  {
  	[]bangumiInfo:array,
  	"msg" :string 
  }
  ```

推荐

- URL: get: "/api/v1/{id}/recommends"

- Req body

  ```json
  nil
  ```

- Rsp body

  ```json
  {
  	[]bangumiInfo:array,
  	"msg": string
  }
  ```

  

#### 信息

##### 查看全部番剧信息

- URL: get: "/api/v1/bangumi/all/{page_id}"

- Req body

  ```json
  nil
  ```

- Rsp body

  ```json
  {
  	[]bangumiInfo:array,
  	"msg": string			
  }
  ```

  

##### 查看单个番剧信息

- URL: get: "/api/v1/bangumi/{b_id} "

- Req body

  ```json
  nil
  ```

- Rsp body

  ```json
  {
  	[]bangumiInfo: array,
  	"msg": string
  }
  ```

  

#### bangumiInfo格式

```json
{
待定
}
```

