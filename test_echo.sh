#!/bin/bash

set -ex



pids=`ps -ux | grep echo_ | grep -v grep  | awk '{print $2}'`
if [[ $pids != "" ]]; then
    kill -9 $pids
fi


cd bin

nohup ./echo_client --server_id=01 > /dev/null 2>&1 &
nohup ./echo_client --server_id=02 > /dev/null 2>&1 &
nohup ./echo_client --server_id=03 > /dev/null 2>&1 &
nohup ./echo_client --server_id=04 > /dev/null 2>&1 &
nohup ./echo_client --server_id=05 > /dev/null 2>&1 &
nohup ./echo_client --server_id=06 > /dev/null 2>&1 &
nohup ./echo_client --server_id=07 > /dev/null 2>&1 &
nohup ./echo_client --server_id=08 > /dev/null 2>&1 &
nohup ./echo_client --server_id=09 > /dev/null 2>&1 &
nohup ./echo_client --server_id=10 > /dev/null 2>&1 &
nohup ./echo_client --server_id=11 > /dev/null 2>&1 &
nohup ./echo_client --server_id=12 > /dev/null 2>&1 &
nohup ./echo_client --server_id=13 > /dev/null 2>&1 &
nohup ./echo_client --server_id=14 > /dev/null 2>&1 &
nohup ./echo_client --server_id=15 > /dev/null 2>&1 &
nohup ./echo_client --server_id=16 > /dev/null 2>&1 &
nohup ./echo_client --server_id=17 > /dev/null 2>&1 &
nohup ./echo_client --server_id=18 > /dev/null 2>&1 &
nohup ./echo_client --server_id=19 > /dev/null 2>&1 &
nohup ./echo_client --server_id=20 > /dev/null 2>&1 &
nohup ./echo_client --server_id=21 > /dev/null 2>&1 &
nohup ./echo_client --server_id=22 > /dev/null 2>&1 &
nohup ./echo_client --server_id=23 > /dev/null 2>&1 &
nohup ./echo_client --server_id=24 > /dev/null 2>&1 &
nohup ./echo_client --server_id=25 > /dev/null 2>&1 &
nohup ./echo_client --server_id=26 > /dev/null 2>&1 &
nohup ./echo_client --server_id=27 > /dev/null 2>&1 &
nohup ./echo_client --server_id=28 > /dev/null 2>&1 &
nohup ./echo_client --server_id=29 > /dev/null 2>&1 &
nohup ./echo_client --server_id=30 > /dev/null 2>&1 &


./echo_server --server_id=1 --log_to_stdout=true

