NAME	= gommon

GO_CC				= go

COVER_FILE	= .coverage.out
COVER_VALUE	= -cover -coverprofile=$(COVER_FILE) -covermode=atomic
TEST_COVER	= $(COVER_VALUE)
TEST_ARGS		= -v -short
TEST_FILTER = #-run <pattern>
app_name = github.com/burgesQ/$(NAME)
TEST_FILES	= ./...

TEST				= $(GO_CC) test $(TEST_COVER) $(TEST_ARGS) $(TEST_FILES) $(TEST_FILTER)
LINT				= golangci-lint run

$(NAME): lint test

## lint: lint the go code
lint: ; $(LINT)

## test: run the go test w/ coverage
test: ; $(TEST)

## tidy: Tidy up the deps
tidy:
	@echo "Tidy up deps ..."
	@cd $(SRC) ; $(GO_CC) mod tidy
	@echo "Deps tidyed up"

.PHONY: help
## help: Prints this help message
 help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
