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
- [protobuf](tools/protoc-gen-vmicro/examples/greeter) ，protobuf 协议生成例子
- [flags](examples/flags) ，命令行参数例子
- [registry](examples/registry) ，服务发现例子
- [echo](examples/echo) ，回显测试

## 基准测试

- [回显测试报告]()



## TODO

- [异步 RPC 广播](doc/异步RPC广播使用界面设计.md)
- Registry 新增插件: consul
- Transport 新增插件： http
- echo 回显测试，性能报告
- 新增吞吐量测试，性能报告
- 更多的例子
- 更多的测试、单元测试


## 贡献

代码很大一部分来源于 [go-micro](https://github.com/micro/go-micro) ，太多，不再一一注释说明。特此声明
