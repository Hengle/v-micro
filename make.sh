#!/bin/bash

set -e

CUR_DIR=$PWD
SRC_DIR=$PWD
cmd=$1

LINT_VER=v0.0.0-20190409202823-959b441ac422
export GOPROXY=https://goproxy.io
export GOPATH=~/go
export GOBIN=~/go/bin

case $cmd in
    lint) $0 dep && $GOBIN/golint -set_exit_status $(go list $SRC_DIR/...) ;;
    test) go test -short $(go list $SRC_DIR/...) ;;
    race) go test -race $(go list $SRC_DIR/...) ;;
    coverage) rm -rf coverage.cov && go test -race $(go list $SRC_DIR/...) -coverprofile=coverage.cov -covermode=atomic ;;
    dep) go get -v golang.org/x/lint@$LINT_VER && cd $GOPATH/pkg/mod/golang.org/x/lint@$LINT_VER/ && go install ./... && cd $CUR_DIR ;;
    build) export GOBIN=$SRC_DIR/bin && go install $SRC_DIR/... ;;
    *) echo 'This script is used to test, build, and publish go code'
       echo ''
       echo 'Usage:'
       echo ''
       echo '        lint               Lint the files'
       echo '        test               Run unittests'
       echo '        race               Run data race detector'
       echo '        coverage           Generate global code coverage report'
       echo '        dep                Get the dependencies'
       echo '        build              Build the binary file'
       echo ''
    ;;
esac

cd $CUR_DIR