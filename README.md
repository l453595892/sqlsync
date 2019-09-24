### sqlsync
> 1.此项目是基于raft协议实现的基础版本  
> 2.目前是一个前期的版本支持本地mysql和pgsql的分布式场景下同步数据库，但并不像etcd一样捆绑底层存储  
> 3.目前的dockerfile是未测试版本，有时间会去测试
---
### 执行项目
1.编译项目
```
go build -o sqlsync cmd/main.go
```
2.设置环境变量（datapath为raft日志存储路径，type为postgres或mysql）
```
export DATAPATH=data1
export CONFIG_TYPE=postgres
export CONFIG_USERNAME=root
export CONFIG_PASSWORD=123456
export CONFIG_HOST=127.0.0.1
export CONFIG_PORT=5433
export CONFIG_DATABASE=root
```
3.运行
```
./sqlsync -h localhost -p 4001
```
```
./sqlsync -h localhost -p 4002 -join localhost:4001
```
---
TODO  
1.优化Docker支持  
2.优化sql支持  
3.优化raft内容和压测