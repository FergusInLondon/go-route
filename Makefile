.PHONY: clean

clean:
	rm -f server
	rm -f client

server:
	go build -o server cmd/server/main.go

client:
	go build -o client cmd/client/main.go

build: clean server client
