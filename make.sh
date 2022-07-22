#!/bin/bash 

set -e 

rm -rf ./bin 2>/dev/null || :
mkdir -pv bin 

GOOS=linux GOARCH=arm GOARM=7 go build  -o ./bin/timeIt-linux-armv7l 
GOOS=linux GOARCH=arm GOARM=6 go build  -o ./bin/timeIt-linux-armv6l 
GOOS=windows GOARCH=amd64 go build -o ./bin/timeIt-win-amd64.exe 
GOOS=linux GOARCH=amd64 go build -o ./bin/timeIt-linux-amd64 

cd bin && for i in timeIt*; do 
  sha256sum $i > $i.sha256
done

