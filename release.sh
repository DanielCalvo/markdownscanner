#!/bin/bash
#Github actions? Gitlab? Tekton? Naaaaaaaaaah

env GOOS=linux GOARCH=arm go build
rsync -av . daniel@192.168.1.112:/home/daniel/markdownscanner
rm markdownscanner