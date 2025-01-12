#!/usr/bin/bash
env GOOS=linux GOARCH=arm64 go build tool.go  -o nshptt_arm64