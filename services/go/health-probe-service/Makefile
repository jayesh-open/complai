.PHONY: build run test lint clean

SERVICE_NAME := health-probe-service
BINARY := bin/$(SERVICE_NAME)

build:
	go build -o $(BINARY) ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./... -v -race

lint:
	go vet ./...

clean:
	rm -rf bin/
