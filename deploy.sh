#!/usr/bin/env bash

#YOLOMAX
#Proper release work (ex: Dockerfile and instructions) is pending.

go build main.go
rm -rf ./repositories
rsync -avz * ubuntu@dcalvo.dev:/home/ubuntu/markdownscanner

#Remotely:
#Manually kill any processes
#./main -slowscan=true &> mdscanner.log &

