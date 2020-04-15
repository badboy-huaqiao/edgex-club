#!/bin/bash

DIR=$PWD
CMD=../cmd/edgex-club
function cleanup {
	pkill edgex-club
}

cd $CMD
exec -a edgex-club ./edgex-club &
cd $DIR

trap cleanup EXIT

while : ; do sleep 1 ; done