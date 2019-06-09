### go框架目录结构

    ├── app                                 //应用目录
    │   ├──http                             //restful接口
    │   │    ├──controller                  //控制层
    │   │    ├──middleware                  //中间件
    │   │    ├──goserver.sh                 //启动脚本
    │   │    ├──main.go                     //主文件
    │   ├──websocket                        //websocket服务
    │   └──rpc                              //rpc服务
    ├── conf
    │   ├── conf.go                         //配置文件对应struct
    │   ├── redisKeys.go                    //redis key文件
    │   ├── retcode.go                      //返回错误码配置文件
    │   └── toml
    │       ├──dev.toml                     //开发环境配置文件 
    │       ├──test.toml                    //测试环境配置文件
    │       └──online.toml                  //线上环境配置文件 
    ├── dao                                 //数据层
    ├── model                               //结构实体层
    ├── service                             //服务层
    ├── lib                                 //公用类库包目录
    │   ├── encrypt                         //数据加解密类库
    │   ├── logger                          //日志类库.
    │   ├── mail                            //mail类库
    │   ├── myredis                         //redis类库
    │   └── util                            //公用函数类库
    ├── vendor                              //包依赖
    ├── go.mod                              //go module文件  
    ├── go.sum                              //go sum文件
    └── README.md                           //框架说明文件

### 安装说明

        本框架依赖go module 
        Requirments Go version>=1.12 or go1.11  GO111MODULE=on
        
#### 框架包名为goframe，如果要创建自己包名项目，请使用如下工具

        go get -u github.com/operahl/pkgreplace 
         
        pkgreplace goframe myframe 
        
        export GO111MODULE=on
                
        go mod vendor

### 启动说明

     go run app/http/main.go ./conf/toml/dev.toml


### 注意事项

    1. app下的每个目录都是独立进程，根据需求创建自己的应用

    2. 路由文件写到main.go里, controller 编写参数验证和结果返回,service写业务逻辑,dao层写底层数据交换,model层写数据结构，定义表结构或者返回数据实体

    3. 调用方式是 main->controller->service->dao