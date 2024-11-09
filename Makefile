.PHONY: all
all: build

.PHONY: build
build:
	go build -o bin/lazynx cmd/lazynx/main.go

.PHONY: run
run:
	go run cmd/lazynx/main.go

.PHONY: build-run
build-run: build
	./bin/lazynx
