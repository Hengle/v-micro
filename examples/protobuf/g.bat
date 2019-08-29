REM download protoc
REM download protoc-gen-vmicro

set GOPROXY=https://goproxy.io
go get github.com/fananchong/protoc-gen-vmicro
go get github.com/golang/protobuf
go list -m -f "{{.Dir}}" github.com/fananchong/protoc-gen-vmicro > temp
set /P DEP1=<temp
go list -m -f "{{.Dir}}" github.com/golang/protobuf > temp
set /P DEP2=<temp
echo %DEP1%
echo %DEP2%

protoc -I. -I%DEP1% -I%DEP2% --vmicro_out=. greeter.proto