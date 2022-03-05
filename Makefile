

.PHONY: setup
setup:
	go mod tidy
	go install github.com/matryer/moq@latest

.PHONY: build
build:
	go mod tidy
	rm -rf dist && mkdir dist
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o dist ./...

docker.start.debug ::
	docker-compose -f docker-compose.debug.yml up -d

docker.stop ::
	docker-compose down