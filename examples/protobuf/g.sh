#!/bin/bash

set -ex

# download protoc
if [ ! -f "./protoc" ]; then
    wget -P ./temp https://github.com/protocolbuffers/protobuf/releases/download/v3.9.1/protoc-3.9.1-linux-x86_64.zip
    unzip -d ./temp ./temp/protoc-3.9.1-linux-x86_64.zip
    cp -f ./temp/bin/protoc .
    rm -rf ./temp
fi

# download protoc-gen-vmicro
if [ ! -f "./protoc-gen-vmicro" ]; then
    wget https://github.com/fananchong/protoc-gen-vmicro/releases/download/v1.0.0/protoc-gen-vmicro
    chmod +x ./protoc-gen-vmicro
fi

export PATH=$PATH:$PWD
export GOPROXY=https://goproxy.io
go get github.com/golang/protobuf
go get github.com/fananchong/protoc-gen-vmicro
DEP1=`go list -m -f "{{.Dir}}" github.com/golang/protobuf`
DEP2=`go list -m -f "{{.Dir}}" github.com/fananchong/protoc-gen-vmicro`
protoc -I. -I$DEP1 -I$DEP2 --vmicro_out=. greeter.proto
