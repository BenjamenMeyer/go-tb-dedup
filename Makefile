COVERAGE_FILE=.cover
BINARY_NAME="go-tb-dedup"
all:
	go build -o $(BINARY_NAME) -race .

test:
	go test -cover -covermode=atomic -coverprofile=$(COVERAGE_FILE) ./... -v

coverage: test
	go tool cover -html=$(COVERAGE_FILE)
	go tool cover -func=$(COVERAGE_FILE)

format:
	go fmt ./...
