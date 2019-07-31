## 使用界面

使用界面基本上抄 go-micro ，但有细微区别，举例如下

##### 1. 协议定义

```protobuf
syntax = "proto3";
package proto;

service Say {
    rpc Hello(Request) returns (Response) {}
}

message Request {
    string name = 1;
}

message Response {
    string msg = 1;
}
```

- 需求点： rpc 调用，在客户端、服务器端都不要阻塞

##### 2. 服务器端代码

```go
type Say struct{}

func (s *Say) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
    log.Print("Received Say.Hello request")
    rsp.Msg = "Hello " + req.Name
    return nil
}

func main() {
    service := micro.NewService(
        micro.Name("greeter"),
    )
    service.Init()

    // Register Handlers
    proto.RegisterSayHandler(service.Server(), new(Say))

    // Run server
    if err := service.Run(); err != nil {
        log.Fatal(err)
    }
}
```
- 需求点： Hello 的 rsp 消息回复不阻塞 OnRecv 协程
- proto.RegisterSayHandler 是自动生成代码
- proto.RegisterSayHandler 可以多次调用，方便组织 proto 文件

##### 3. 客户端代码

```go
type Say struct{}

func (s *Say) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
    log.Print("Received Say.Hello Response")
    return nil
}

func main() {
    service := micro.NewService()
    service.Init()

    // Use the generated client stub
    cl := proto.NewSayService("greeter", new(Say), service.Client())

    // Make request
    cl.Hello(context.Background(), &proto.Request{
        Name: "John",
    })
}
```

客户端使用界面与 go-micro 区别比较大：
- 需求点： proto.NewSayService 自动生成代码， 多接收 new(Say) 消息回调处理
- 需求点： cl.Hello 不再阻塞，发送调用就返回
- proto.NewSayService 可以多次调用，方便组织 proto 文件

## 实现上与 go-micro 不同的地方

- 上述使用界面上的不同
- server / client 底层节点间只维持 1 个连接，不是 call 1 次建立一个连接或有 pool 这种
