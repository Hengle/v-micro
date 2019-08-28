## 使用界面


##### 1. 协议定义

```protobuf
syntax = "proto3";
import "micro/broadcast.proto";
package proto;

service Say {
    rpc Hello(Request) returns (micro.NoReply) { option (micro.Broadcast) = true; }
}

message Request {
    string name = 1;
}
```

- 必须 import broadcast.proto ，v-micro 扩展定义了 method option ,来标记一个方法是用来广播的
- micro.Noreply 非必须，使用公共的，则减少一个消息定义
- 标记为 option (micro.Broadcast) = true; 的 RPC 方法, 服务器处理完毕不会回消息
- 标记为 option (micro.Broadcast) = true; 的 RPC 方法, 客户端广播完毕不会触发消息回调

##### 2. 服务器端代码

同`异步 RPC 调用`，代码展示略


##### 3. 客户端代码

```go
func main() {
    service := micro.NewService()
    service.Init()

    // Use the generated client stub
    cl := proto.NewSayService("greeter", nil, service.Client())

    // Broadcast request for all [service greeter]
    cl.BroadcastHello(context.Background(), &proto.Request{
        Name: "John",
    })

    // Broadcast request with filter in [service greeter]
    cl.BroadcastHello(context.Background(), &proto.Request{
        Name: "John",
    }, Filter("latest"))
}

// Filter will filter the version of the service ( for example)
func Filter(v string) client.CallOption {
	filter := func(services []*registry.Service) []*registry.Service {
		var filtered []*registry.Service

		for _, service := range services {
			if service.Version == v {
				filtered = append(filtered, service)
			}
		}

		return filtered
	}

	return client.WithSelectOption(selector.WithFilter(filter))
}
```

- NewSayService 、 BroadcastHello 方法通过 protoc_gen_vmicro 自动生成
- Filter 方法仅为例子演示
- 给所有不同类型的微服务广播消息，暂时不做
