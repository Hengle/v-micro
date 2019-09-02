# v-micro
[![Build Status](https://www.travis-ci.org/fananchong/v-micro.svg?branch=master)](https://www.travis-ci.org/fananchong/v-micro) [![codecov](https://codecov.io/gh/fananchong/v-micro/branch/master/graph/badge.svg)](https://codecov.io/gh/fananchong/v-micro)

可插拔的微服务框架（参考 go-micro）

## 起源

起源及开 v-micro 的原因，见 [doc/intro.md](doc/intro.md)

## 改进

- [异步 RPC 调用](doc/异步RPC调用使用界面设计.md)
- [异步 RPC 广播](doc/异步RPC广播使用界面设计.md)

## 依赖

- go1.12 +

## 编译

- Windows 编译
  ```shell
  make build
  ```

- Linux 编译
  ```shell
  ./make.sh build
  ```

## 例子

- [hello](examples/hello) ，入门例子
- [protobuf](examples/protobuf) ，protobuf 协议生成例子
- [flags](examples/flags) ，命令行参数例子
- [registry](examples/registry) ，服务发现例子
- [echo](examples/echo) ，回显测试
- [broadcast](examples/broadcast) ，广播消息例子
- [metadata](examples/metadata) ，客户端元数据、服务器端元数据例子
- [filter](examples/filter) ，过滤器例子
- [wrapper](examples/wrapper) ，客户端、服务器端包装器例子


## 基准测试

- [回显测试报告](examples/echo/README.md)

## TODO

- Registry 新增插件： consul
- Transport 新增插件： http
- 吞吐量测试，性能报告
- 更多的例子
- 更多的测试、单元测试
- 回显测试性能优化
  - 客户端回调中的 Request 对象目前由服务器传给客户端，多消耗网络流量与 CPU 消耗
  - RPC 过程中减少分配内存，参考官方 RPC 实现
  - 精简每次通信中的默认元数据，绝大部分不会被用到，白白暂流量


## 贡献

代码很大一部分来源于 [go-micro](https://github.com/micro/go-micro) ，太多，不再一一注释说明。特此声明
