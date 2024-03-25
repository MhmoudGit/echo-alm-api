.PHONY: run

run:
	@swag fmt && swag init && go run .
