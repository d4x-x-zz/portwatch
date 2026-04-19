.PHONY: build test lint run clean

BIN := portwatch
CMD := ./cmd/portwatch

build:
	go build -o $(BIN) $(CMD)

test:
	go test ./...

lint:
	golangci-lint run ./...

run: build
	./$(BIN) -config portwatch.toml

clean:
	rm -f $(BIN)
