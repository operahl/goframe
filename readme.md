### 安装说明

        本框架依赖go module 
        Requirments Go version>=1.12 or go1.11  GO111MODULE=on
        
### arc 配置

1 设定`phabricator.uri`属性
```shell
arc set-config phabricator.uri https://pha.o-pay.in
```

2 设定 `vim`为编辑器
```shell
arc set-config editor /usr/bin/vim
```

3 安装证书，运行arc，会提示你用浏览器打开一个链接，获取一个Token，然后粘贴获得的Token按回车即可。
```shell
arc install-certificate
```


### 注意事项

    1. app下的每个目录都是独立进程，根据需求创建自己的应用

    2. 路由文件写到main.go里, controller 编写参数验证和结果返回,service写业务逻辑,dao层写底层数据交换,model层写数据结构，定义表结构或者返回数据实体
