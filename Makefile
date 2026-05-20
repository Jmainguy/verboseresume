.PHONY: all test build build-prod run clean docker fmt ci lint

BINARY ?= verboseresume
# Legacy local name many folks still run after `make build`.
LEGACY_BINARY ?= resumeGen
IMAGE ?= zot.soh.re/verboseresume/verboseresume

all: test build

# Mirrors .github/workflows/golang-ci.yml (install golangci-lint for lint).
ci: fmt
	go vet ./...
	$(MAKE) test
	$(MAKE) build-prod

fmt:
	@test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)

lint:
	golangci-lint run --timeout=5m

test:
	go test ./...

build:
	go build -o $(BINARY) .

build-prod:
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(BINARY) .
	@cp -f $(BINARY) $(LEGACY_BINARY)

run: build-prod
	./$(LEGACY_BINARY)

docker:
	docker build -t $(IMAGE):local .

clean:
	rm -f $(BINARY)
