#!/bin/bash
cd /root/3d-library
export PATH=/usr/local/go/bin:$PATH

echo "Starting worker..."
go run cmd/worker/main.go > /tmp/worker.log 2>&1 &
echo "Worker PID: $!"

sleep 2

echo "Starting web server..."
go run cmd/web/main.go > /tmp/server.log 2>&1 &
echo "Server PID: $!"

echo
echo "✓ Services started"
echo "  Server: http://192.168.3.26:3000"
echo "  Worker log: tail -f /tmp/worker.log"
echo "  Server log: tail -f /tmp/server.log"
