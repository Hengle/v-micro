#!/bin/bash

set -e

CUR_DIR=$PWD
SRC_DIR=$PWD
cmd=$1

export GOPROXY=https://goproxy.io

cd $SRC_DIR

case $cmd in
    coverage) go test -race ./... -coverprofile=coverage.cov -covermode=atomic ;;
    build) export GOBIN=$SRC_DIR/bin && go install ./... ;;
    *) echo 'This script is used to test, build, and publish go code'
       echo ''
       echo 'Usage:'
       echo ''
       echo '        coverage           Generate global code coverage report'
       echo '        build              Build the binary file'
       echo ''
    ;;
esac

cd $CUR_DIR