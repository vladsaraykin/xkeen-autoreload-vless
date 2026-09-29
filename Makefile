APP=xkeen-autoreload-vless
VERSION=$(shell tr -d '[:space:]' < VERSION)

.PHONY: test build build-mipsle fmt vet
fmt:
	gofmt -w ./cmd ./internal

test:
	go test ./...

vet:
	go vet ./...

build:
	mkdir -p dist
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" -o dist/$(APP) ./cmd/xkeen-autoreload

build-mipsle:
	mkdir -p dist
	CGO_ENABLED=0 GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" -o dist/$(APP)-linux-mipsle ./cmd/xkeen-autoreload
