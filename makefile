.PHONY: run

run:
	@golangci-lint run && swag fmt && swag init && go run .
