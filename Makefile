.PHONY: test vet fmt build docker all

test:
	go test ./... -count=1 -race -cover

vet:
	go vet ./...

fmt:
	gofmt -w .

fmt-check:
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then echo "$$unformatted"; exit 1; fi

build:
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/zoho-cliq-release-notifier ./cmd/zoho-cliq-release-notifier

docker:
	docker build -t zoho-cliq-release-notifier:local .

all: fmt-check vet test build
