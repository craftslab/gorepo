#!/bin/bash

go env -w GOPROXY=https://goproxy.cn,direct

CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/gorepo main.go
CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -ldflags="-s -w" -o bin/gorepo.exe main.go

apt install upx

upx bin/gorepo
upx bin/gorepo.exe
