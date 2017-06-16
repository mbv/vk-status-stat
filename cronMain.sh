#!/usr/bin/env bash
cd /var/work/vk-status-stat/
export GOROOT=/usr/local/go
export GOPATH=/var/work/go_libs:/var/work/vk-status-stat/
export PATH=$PATH:$GOROOT/bin
go run src/main.go 2> logFatal.txt