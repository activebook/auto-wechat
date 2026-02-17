.PHONY: build run clean test fmt tidy

BINARY_NAME=auto-wechat

build:
	go build -o $(BINARY_NAME) main.go

run: build
	./$(BINARY_NAME)

clean:
	go clean
	rm -f $(BINARY_NAME)

test:
	go test ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy
