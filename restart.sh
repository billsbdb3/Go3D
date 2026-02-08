#!/bin/bash
pkill -9 -f "go run"
pkill -9 -f "3d-library"
sleep 2
cd /root/3d-library
export PATH=/usr/local/go/bin:$PATH
nohup go run cmd/web/main.go > /tmp/server.log 2>&1 &
sleep 3
tail -5 /tmp/server.log
