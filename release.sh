#!/bin/bash
#Just throwing it over the fence yolo

env GOOS=linux go build
ssh daniel@192.168.1.112 'mkdir -p /home/daniel/Projects/markdownscanner' #locally installed rsync is old and doesn't have `--mkpath` available
rsync -av --no-perms --chown=daniel:daniel . daniel@192.168.1.112:/home/daniel/Projects/markdownscanner
rm markdownscanner