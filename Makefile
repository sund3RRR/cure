.PHONY: run build deploy lint test

run:
	go run cmd/cure/main.go

build:
	mkdir -p bin/
	cd cmd/cure && go build -o ../../bin/cure

deploy:
	docker build -t cure-builder -f deploy/Dockerfile --output=deploy/bin .

test:
	WORKSPACE_DIR=$(shell pwd) go test -v -cover ./...
lint:
	golangci-lint run
