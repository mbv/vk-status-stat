#!/usr/bin/env bash
cd /app
export GOPATH=/app:/go
go run src/main.go 2> logFatal.txt