build:
	@go build -o bin/retoeNFA

run: build
	@./bin/retoeNFA

test:
	@go test -v ./...