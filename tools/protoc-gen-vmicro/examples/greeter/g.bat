set WORK_DIR=%~dp0

cd ..\..\..\..\
go install ./...
cp /y bin\protoc-gen-vmicro.exe %WORK_DIR%

cd %WORK_DIR%
protoc --proto_path=. --vmicro_out=. greeter.proto