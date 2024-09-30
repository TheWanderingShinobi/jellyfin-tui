# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=jellyfin-tui
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/jellyfin-tui

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/jellyfin-tui
	./$(BINARY_NAME)

deps:
	$(GOGET) -v -d ./...

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/jellyfin-tui

docker-build:
	docker build -t $(BINARY_NAME):latest .

.PHONY: all build test clean run deps build-linux docker-build
