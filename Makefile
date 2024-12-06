.PHONY: all
all: build

.PHONY: build
build:
	go build ./cmd/clipy

install:
	go install ./cmd/clipy

lint:
	golangci-lint run

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -f clipy debug.log
