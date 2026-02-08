.PHONY: dev build run migrate-up migrate-down

dev:
	cd cmd/web && /usr/local/go/bin/go run main.go

build:
	/usr/local/go/bin/go build -o bin/web cmd/web/main.go
	/usr/local/go/bin/go build -o bin/worker cmd/worker/main.go

run:
	./bin/web

deps:
	/usr/local/go/bin/go mod download

test:
	/usr/local/go/bin/go test ./...

clean:
	rm -rf bin/
