.PHONY: all
all: vet test build

.PHONY: build
build:
	go build  ./cmd/poslog

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	rm -f poslog
