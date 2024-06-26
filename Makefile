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
	@golangci-lint run
