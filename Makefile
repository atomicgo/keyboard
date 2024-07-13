##################################################################
#                          IMPORTANT NOTE                        #
#                                                                #
# This file is synced with https://github.com/atomicgo/template  #
# Please make changes there.                                     #
##################################################################

test:
	@echo "# Running tests..."
	@go test -v ./...

lint:
	@echo "# Linting..."
	@echo "## Go mod tidy..."
	@go mod tidy
	@echo "## Fixing whitespaces..."
	@wsl --allow-cuddle-declarations --force-err-cuddling --force-case-trailing-whitespace 3 --fix ./...
	@echo "## Running golangci-lint..."
	@golangci-lint run
