test:
	@echo "# Running tests..."
	@go test -v ./...

lint:
	@echo "# Linting..."
	@golangci-lint run
