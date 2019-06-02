
all: lint test
	go build

lint:
	golangci-lint run --deadline 300s

test:
	go test ./...
