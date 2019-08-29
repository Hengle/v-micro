## 依赖

需要以下 2 进制程序：

- protoc
- protoc_gen_vmicro

生成、获取方法，参见：[examples/protobuf/README.md](../../../examples/protobuf/README.md)

可以拷贝至本目录或放到系统目录


## 生成命令

```shell
protoc -I. --vmicro_out=. echo.proto
```

更详细的，参见：[examples/protobuf/README.md](../../../examples/protobuf/README.md)