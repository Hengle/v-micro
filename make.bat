chcp 65001

set CUR_DIR=%~dp0
set SRC_DIR=%~dp0
set CMD="%1"

set GOPROXY=https://goproxy.io
set GOBIN=%CUR_DIR%bin

if %CMD%=="build" (
    gofmt -w -s .
    go build ./...
)^
else if %CMD%=="race" (
    go test -race ./...
)

cd %CUR_DIR%