
.PHONY: setup
setup: setup/go

.PHONY: setup/go
setup/go:
	go mod tidy
	go install github.com/matryer/moq@latest