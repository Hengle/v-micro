chcp 65001

set CUR_DIR=%~dp0
set SRC_DIR=%~dp0
set CMD="%1"

set GOPROXY=https://goproxy.io
set GOBIN=%CUR_DIR%bin

if %CMD%=="build" (
    gofmt -w -s .
    go install ./...
)^
else if %CMD%=="race" (
    go test -race ./...
)^
else if %CMD%=="test" (
    go test -v ./...
)^
else if %CMD%=="start" (
    cd %GOBIN%
    start server.exe --server_id=1 --log_to_stdout=true
    start server.exe --server_id=2 --log_to_stdout=true
    start server.exe --server_id=3 --log_to_stdout=true
)


cd %CUR_DIR%