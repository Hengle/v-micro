## protoc

可以从 https://github.com/protocolbuffers/protobuf/releases 页面下载

比如:

- Windows
  - https://github.com/protocolbuffers/protobuf/releases/download/v3.9.1/protoc-3.9.1-win64.zip

- Linux
  - https://github.com/protocolbuffers/protobuf/releases/download/v3.9.1/protoc-3.9.1-linux-x86_64.zip

## protoc-gen-vmicro

- Windows
  - https://github.com/fananchong/protoc-gen-vmicro/releases/download/v1.0.0/protoc-gen-vmicro.exe

- Linux
  - https://github.com/fananchong/protoc-gen-vmicro/releases/download/v1.0.0/protoc-gen-vmicro


## 生成 pb.go 文件

比如，拷贝 protoc 、 protoc-gen-vmicro 至本目录下。执行：

- Windows
    ```shell
    protoc -I. --vmicro_out=. greeter.proto
    ```

- Linux
    ```shell
    export PATH=$PATH:$PWD
    protoc -I. --vmicro_out=. greeter.proto
    ```

即可生成 greeter.pb.go


如果带有 broadcast 选项的，参考 g.bat 、 g.sh
