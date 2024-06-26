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
	@echo "## Fixing whitespaces..."
	@wsl --allow-cuddle-declarations --force-err-cuddling --fix ./...
	@echo "## Running golangci-lint..."
	@golangci-lint run
