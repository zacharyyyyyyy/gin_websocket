# gin_websocket
## 简易客服系统

若为windows 需先执行
``````
set GOOS=linux
``````
编译
````
go build main.go 
````
构建镜像(需要有alpine:latest镜像)
````
docker build -t go_server .
````
运行容器
````
docker-compose up -d
````
#### 使用流程
后台
* 调用 /admin/login登录
* 调用 /admin/ws 注册当前管理员，接受websocket消息
* 调用 /admin/ws/info 获取当前管理员下用户信息
* 调用 /admin/ws/link 传入info接口分配的ws_key即可进行链接
* 调用 /admin/ws/close 取消注册，即不在分配websocket链接给当前管理员

api
* 调用 /api/link  进行链接

websocket传参格式
```
type目前为ping与chat两种
```
```json
{
  "content": "test",
  "type": "chat"  
}
```


 


###文件结构
<details>
<pre><code>
├── main.go
│
├── conf  //配置文件目录
│   └── sample-config.ini
│
├── controller
│   ├── admin //后台相关接口
│   ├── api  //外部相关接口
│   ├── perf //perf相关接口
│   └── base_controller.go //controller统一方法 目前用于统一返回信息
│
├── dao
│
├── lib //
│   ├── config  //配置读取相关
│   │   ├── config.go  //统一配置入口
│   │   ├── model.go  //mysql相关配置
│   │   ├── kafka.go  //kafka相关配置
│   │   ├── mq.go  //rabbitmq相关配置
│   │   ├── redis.go  //redis相关配置
│   │   └── websocket.go  //websocket相关配置
│   ├── kafka  //kafka相关函数
│   ├── logger  //zaplogger相关函数
│   ├── mq  //rabbitmq相关函数
│   ├── redis  //redis相关函数
│   ├── session  //session相关函数
│   ├── tools //常用工具函数
│   ├── validator //controller层参数校验相关函数
│   │   ├── validator.go  //校验统一入口
│   │   └── dao_validator.go  //需要dao层的相关校验
│
│   ├── log  //日志
│   │   ├── api  //api相关日志
│   │   ├── model  //model相关日志
│   │   ├── runtime  //运行产生的其他日志
│   │   ├── service  //service层相关日志
│   │   └── taskqueue  //延时队列相关日志
│
│   ├── middleware  //中间件
│   │   ├── global_middleware  //全局中间件
│   │   │   ├── cors.go  //跨域相关
│   │   │   ├── http_recover.go //外层recover
│   │   │   ├── http_trace.go  //请求追踪并记录
│   │   │   └── no_route.go   //请求不存在时逻辑处理
│   │   └── router_middleware  //路由中间件
│   │   │   ├── auth.go   //权限
│   │   │   └── login_limit.go  //登录错误次数限制
│
│   ├── model  //模型层
│   │   └── base_model.go  //模型统一入口
│
│   ├── router  //路由
│   │   ├── admin.go  //后台相关路由
│   │   ├── api.go  //api相关路由
│   │   └── router.go  //router统一入口
│
│   ├── service  //service层
│   │   ├── admin  //后台相关逻辑
│   │   ├── taskqueue  //延时队列相关逻辑(依赖数据库)
│   │   │   ├── task  //延时队列消费
│   │   │   │   ├── task.go  //统一逻辑入口
│   │   │   │   ├── mq_consumer.go  //rabbitmq超时重发相关逻辑
│   │   │   │   └── kafka_consumer.go  //kafka重发相关逻辑
│   │   │   └── tasakqueue.go 
│   │   ├── tracer  //http_trace 请求记录相关逻辑
│   │   ├── websocket //websocket 相关逻辑
│   │   └── service.go  //service统一入口
│
└── service
</code></pre>
</details>