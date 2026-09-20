.PHONY: all build test clean install

BINARY_NAME=tg
BUILD_DIR=bin

all: build

build:
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) .

test:
	go test -v ./...

install:
	go install .

clean:
	rm -rf $(BUILD_DIR)
