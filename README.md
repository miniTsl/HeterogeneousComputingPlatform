# HeterogeneousComputing


## 架构简介



## 系统行为
1. 用户发送指令让服务器执行
2. 用户发送指令让服务器转交给设备执行
3. 用户发送数据让服务器转交给设备存储

protoc --go_out=./ .\message.proto
protoc --go_out=./ .\register.proto
protoc --go-grpc_out=. .\register.proto

protoc --go-grpc_out=. .\register.proto
protoc --go-grpc_out=. .\profile.proto
## TODO
1. 实现跳板机和一级设备分离
2. 实现一级设备自动识别,热插拔二级设备
3. 实现客户端发送文件到一级设备
4. 为hcp平台添加web入口

### 
1. [cobra使用教程](https://xcbeyond.cn/blog/golang/cobra-quick-start/)
2. [logrus第三方日志库](https://github.com/sirupsen/logrus)