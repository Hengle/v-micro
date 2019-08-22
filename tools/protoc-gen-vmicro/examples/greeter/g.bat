set WORK_DIR=%~dp0
set SRC_DIR=%WORK_DIR%\..\..\..\..\
cd %SRC_DIR%
set GOBIN=%SRC_DIR%bin
go install ./...
copy /y bin\protoc-gen-vmicro.exe %WORK_DIR%

cd %WORK_DIR%
protoc --proto_path=. --vmicro_out=. greeter.proto