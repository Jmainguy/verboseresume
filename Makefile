.PHONY: all test build build-prod run clean docker

BINARY ?= verboseresume
IMAGE ?= ghcr.io/jmainguy/verboseresume

all: test build

test:
	go test ./...

build:
	go build -o $(BINARY) .

build-prod:
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(BINARY) .

run: build
	./$(BINARY)

docker:
	docker build -t $(IMAGE):local .

clean:
	rm -f $(BINARY)
