## protoc

可以从 https://github.com/protocolbuffers/protobuf/releases 页面下载

比如:

- Windows
  - https://github.com/protocolbuffers/protobuf/releases/download/v3.9.1/protoc-3.9.1-win64.zip

- Linux
  - https://github.com/protocolbuffers/protobuf/releases/download/v3.9.1/protoc-3.9.1-linux-x86_64.zip

## protoc-gen-vmicro

- Windows
  ```shell
  git clone https://github.com/fananchong/v-micro.git
  cd v-micro
  make.bat build
  ```

- Linux
  ```shell
  git clone https://github.com/fananchong/v-micro.git
  cd v-micro
  make.sh build
  ```

protoc-gen-vmicro 在 bin 目录下


## 生成 pb.go 文件

比如，拷贝 protoc 、 protoc-gen-vmicro 至本目录下。执行：

```shell
protoc --proto_path=. --vmicro_out=. greeter.proto
```

即可生成 greeter.pb.go
