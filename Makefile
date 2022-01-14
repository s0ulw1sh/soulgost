
.PHONY: all
all:
	go build -o ./bin/soulgost

.PHONY: test
test:
	go test ./...