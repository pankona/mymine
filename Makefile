
all: lint test
	go build

lint:
	gometalinter.v2 --vendor ./...

test:
	go test ./...
