## 依赖

需要以下 2 进制程序：

- protoc
- protoc_gen_vmicro

生成、获取方法，参见：[tools/protoc-gen-vmicro/examples/greeter/README.md](../../../tools/protoc-gen-vmicro/examples/greeter/README.md)

可以拷贝至本目录或放到系统目录


## 生成命令

```shell
protoc --proto_path=. --vmicro_out=. echo.proto
```
