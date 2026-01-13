.PHONY: dev start build test clean

dev:
	go run main.go

start:
	./ecrypto

build:
	go build -o ecrypto

test:
	go test ./...

clean:
	rm -f ecrypto