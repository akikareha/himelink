all: build

build:
	go build -o himelink ./cmd/himelink

clean:
	rm -f himelink

run:
	go run ./cmd/himelink

fmt:
	go fmt ./...

test:
	go test ./...
