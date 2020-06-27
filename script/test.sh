#!/bin/bash

list="gitiles,manifest,repo"

go env -w GOPROXY=https://goproxy.cn,direct

old=$IFS IFS=$','
for item in $list; do
  go test -cover -v $item/*.go
done
IFS=$old
