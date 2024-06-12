#!/bin/bash
#
#
echo usage "$0 [worker] [num]"
test ! -z "$1" && WRK="$1" || WRK=4
test ! -z "$2" && NUM="$2" || NUM=1024

test ! -e "testShell.sh" && echo "error testShell.sh not found" && exit 2

i=1;
while [ $i -le $WRK ]; do
 ./testShell.sh all $NUM &
 let i="i+1"
done
