#!/bin/bash

set -e

CUR_DIR=$PWD
SRC_DIR=$PWD
cmd=$1

LINT_VER=v0.0.0-20190409202823-959b441ac422
export GOPROXY=https://goproxy.io
export GOPATH=~/go
export GOBIN=~/go/bin

# Code coverage generation
function coverage_test(){
    COVERAGE_DIR="${COVERAGE_DIR:-coverage}"
    mkdir -p "$COVERAGE_DIR";
    for package in $(go list $SRC_DIR/...); do
        go test -covermode=count -coverprofile "${COVERAGE_DIR}/${package##*/}.cov" "$package" ;
    done ;
    echo 'mode: count' > "${COVERAGE_DIR}"/coverage.cov ;
    tail -q -n +2 "${COVERAGE_DIR}"/*.cov >> "${COVERAGE_DIR}"/coverage.cov ;
    go tool cover -func="${COVERAGE_DIR}"/coverage.cov ;
    if [ "$1" == "html" ]; then
        go tool cover -html="${COVERAGE_DIR}"/coverage.cov -o coverage.html ;
    fi
    rm -rf "$COVERAGE_DIR";
}

case $cmd in
    lint) $0 dep && $GOBIN/golint -set_exit_status $(go list $SRC_DIR/...) ;;
    test) go test -short $(go list $SRC_DIR/...) ;;
    race) go test -race -short $(go list $SRC_DIR/...) ;;
    coverage) coverage_test ;;
    coverhtml) coverage_test html ;;
    dep) go get -v golang.org/x/lint@$LINT_VER && cd $GOPATH/pkg/mod/golang.org/x/lint@$LINT_VER/ && go install ./... && cd $CUR_DIR ;;
    build) go build $SRC_DIR/... ;;
    *) echo 'This script is used to test, build, and publish go code'
       echo ''
       echo 'Usage:'
       echo ''
       echo '        lint               Lint the files'
       echo '        test               Run unittests'
       echo '        race               Run data race detector'
       echo '        coverage           Generate global code coverage report'
       echo '        coverhtml          Generate global code coverage report in HTML'
       echo '        dep                Get the dependencies'
       echo '        build              Build the binary file'
       echo ''
    ;;
esac

cd $CUR_DIR

