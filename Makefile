# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOLIST=$(GOCMD) list
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet


all: unit build


.PHONY: unit
unit: ## @testing Run the unit tests 
	$(GOFMT) ./...
	$(GOVET) ./cmpp/...
	$(GOTEST) -race -coverprofile=coverage.txt -covermode=atomic $(shell go list ./cmpp/...)

.PHONY: build
build:
	$(GOBUILD) -o ./bin/receiver ./cmd/receiver
	$(GOBUILD) -o ./bin/mockserver ./cmd/mockserver
	$(GOBUILD) -o ./bin/transmitter ./cmd/transmitter


.PHONY: build_linux
build_linux: clean
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -o ./bin/receiver ./cmd/receiver
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -o ./bin/mockserver ./cmd/mockserver
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -o ./bin/transmitter ./cmd/transmitter

.PHONY: clean
clean:
	rm -rf ./bin/ coverage.txt
