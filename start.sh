#!/bin/bash
cd /root/3d-library
export PATH=/usr/local/go/bin:$PATH
go run cmd/web/main.go
