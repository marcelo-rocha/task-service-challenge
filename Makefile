

.PHONY: setup
setup:
	go mod tidy
	go install github.com/matryer/moq@latest

.PHONY: build
build:
	go mod tidy
	rm -rf dist && mkdir dist
	go build -o dist ./...