build:
	@go build -o bin/go_dog

run: build
	@./bin/go_dog

test:
	@go test -v ./...
