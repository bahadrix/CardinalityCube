BIN_DIR := $(GOPATH)/bin
# VERSION := $(shell git describe --tags)
VERSION := "0.1.0"
BUILD := $(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

install-deps:
	@echo "Getting dependencies"
	@go get -v -t ./...
test: install-deps
	go test -v ./...
lint:
	@golint ./...
clean:
	@rm -rf build
	@go clean
build: install-deps
	@mkdir -p build
	@echo "Building server"
	@go build -o build/cubeserver github.com/bahadrix/cardinalitycube/server/service
	@echo "Building client"
	@go build -o build/cubecli github.com/bahadrix/cardinalitycube/cubeclient
	@echo "Done"
