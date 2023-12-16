build:
	go build .

lint:
	golangci-lint run -v .

test:
	go test -v .
