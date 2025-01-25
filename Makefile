.PHONY: run build deploy lint

run:
	go run cmd/cure/main.go

build:
	mkdir -p bin/
	cd cmd/cure && go build -o ../../bin/cure

deploy:
	docker build -t cure-builder -f deploy/Dockerfile --output=deploy/bin .

lint:
	golangci-lint run
