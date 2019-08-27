#!/bin/bash

set -ex

export WORK_DIR=$PWD
export GOPROXY=https://goproxy.io
export GOBIN=$WORK_DIR
go get github.com/golang/protobuf
go get github.com/fananchong/protoc-gen-vmicro
DEP1=`go list -m -f "{{.Dir}}" github.com/golang/protobuf`
DEP2=`go list -m -f "{{.Dir}}" github.com/fananchong/protoc-gen-vmicro`
echo $DEP1
echo $DEP2

if [ ! -f "./protoc" ]; then
    wget -P ./temp https://github.com/protocolbuffers/protobuf/releases/download/v3.9.1/protoc-3.9.1-linux-x86_64.zip
    unzip -d ./temp ./temp/protoc-3.9.1-linux-x86_64.zip
    cp -f ./temp/bin/protoc .
    rm -rf ./temp
fi


export PATH=$PATH:$WORK_DIR
protoc -I. -I$DEP1 -I$DEP2 --vmicro_out=. greeter.proto

