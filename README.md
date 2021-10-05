# go-plugin-bidirectional
go插件双向通讯例子



1.定义grpc 接口
2.定义插件暴露接口
3.定义插件client和server结构体 
    client实现暴露接口的
    server实现grpc的server接口
4.创建插件接口插件结构体 实现server和client
5.实现插件
6.实现条用者
