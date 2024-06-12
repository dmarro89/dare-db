#!/bin/bash
echo "usage: $0 [all|set|get|del] limit"
limit=10
test ! -z "$2" && test "$2" -gt 0 && limit="$2"

start=$(date +%s)

if [[ "$1" = "set" || "$1" = "all" ]]; then
## looping SET
i=1;
while [ $i -le $limit ]; do
 reply=`curl -sd "{\"myKey${i}\":\"Value${i}\"}" "http://localhost:2420/set"`
 echo "set $i: $reply"
 let i="i+1"
done;
fi

if [[ "$1" = "get" || "$1" = "all" ]]; then
## looping GET
i=1;
while [ $i -le $limit ]; do
 echo "get $i: `curl -s "http://localhost:2420/get/myKey${i}"`"
 let i="i+1"
done;
fi


if [[ "$1" = "del" || "$1" = "all" ]]; then
## looping DEL
i=1;
while [ $i -le $limit ]; do
 echo "del $i: `curl -s -X DELETE "http://localhost:2420/delete/myKey${i}"`"
 let i="i+1"
done;
fi

now=$(date +%s)
let took="now-start"
echo "took $took sec"

test ! -z "$1" && exit 0

echo -e "\nget 1-4"
echo -e "get 1: `curl -sX GET "http://localhost:2420/get/myKey1"`"
echo -e "get 2: `curl -sX GET "http://localhost:2420/get/myKey2"`"
echo -e "get 3: `curl -sX GET "http://localhost:2420/get/myKey3"`"
echo -e "get 4: `curl -sX GET "http://localhost:2420/get/myKey4"`"

echo -e "\nset 1-3"
echo -e "set 1: `curl -sX POST -d '{"myKey1":"myValue1"}' "http://localhost:2420/set"`"
echo -e "set 2: `curl -sX POST -d '{"myKey2":"myValue2"}' "http://localhost:2420/set"`"
echo -e "set 3: `curl -sX POST -d '{"myKey3":"myValue3"}' "http://localhost:2420/set"`"

echo -e "\nget 1-4"
echo -e "get 1: `curl -sX GET "http://localhost:2420/get/myKey1"`"
echo -e "get 2: `curl -sX GET "http://localhost:2420/get/myKey2"`"
echo -e "get 3: `curl -sX GET "http://localhost:2420/get/myKey3"`"
echo -e "get 4: `curl -sX GET "http://localhost:2420/get/myKey4"`"

echo -e "\nset NEW value 1-3"
echo -e "new 1: `curl -X POST -d '{"myKey1":"NEWvalue1"}' "http://localhost:2420/set"`"
echo -e "new 2: `curl -X POST -d '{"myKey2":"NEWvalue2"}' "http://localhost:2420/set"`"
echo -e "new 3: `curl -X POST -d '{"myKey3":"NEWvalue3"}' "http://localhost:2420/set"`"

echo -e "\nget 1-4"
echo -e "get 1: `curl -sX GET "http://localhost:2420/get/myKey1"`"
echo -e "get 2: `curl -sX GET "http://localhost:2420/get/myKey2"`"
echo -e "get 3: `curl -sX GET "http://localhost:2420/get/myKey3"`"
echo -e "get 4: `curl -sX GET "http://localhost:2420/get/myKey4"`"

echo -e "\ndel 1-4"
echo -e "del 1: `curl -sX DELETE "http://localhost:2420/delete/myKey1"`"
echo -e "del 2: `curl -sX DELETE "http://localhost:2420/delete/myKey2"`"
echo -e "del 3: `curl -sX DELETE "http://localhost:2420/delete/myKey3"`"
echo -e "del 4: `curl -sX DELETE "http://localhost:2420/delete/myKey4"`"

echo -e "\nget 1-4"
echo -e "get 1: `curl -sX GET "http://localhost:2420/get/myKey1"`"
echo -e "get 2: `curl -sX GET "http://localhost:2420/get/myKey2"`"
echo -e "get 3: `curl -sX GET "http://localhost:2420/get/myKey3"`"
echo -e "get 4: `curl -sX GET "http://localhost:2420/get/myKey4"`"




