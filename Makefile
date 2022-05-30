.PHONY: all
all: vet build

.PHONY: build
build:
	go build ./cmd/poslog

.PHONY: vet
vet:
	go vet ./...

.PHONY: clean
clean:
	rm -f poslog poslog.exe
