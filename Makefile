.PHONY: client
.PHONY: server

all: vet client server

vet:
	@echo "Checking for code issues"
	go vet ./...

client:
	@echo "Removing the client binary"
	rm -f bin/client
	@echo "Building the client binary"
	go build -o bin/client cmd/client/main.go

server:
	@echo "Removing the server binary"
	rm -f bin/server
	@echo "Building the server binary"
	go build -o bin/server cmd/server/main.go