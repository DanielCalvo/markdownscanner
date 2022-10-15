#!/bin/bash
#Github actions? Gitlab? Tekton? Naaaaaaaaaah

env GOOS=linux GOARCH=arm go build
rsync -av --no-perms --chown=root:root . daniel@192.168.31.131:/disk/markdownscanner #Has to be root for now, vfat is terrible
rm markdownscanner