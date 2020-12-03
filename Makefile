.PHONY: client
.PHONY: server

all: yarn vet client server node

yarn:
	@echo "Install Node Modules"
	yarn --cwd notifier/
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
node:
	@echo "Running node server on port 9000"
	node notifier/notifier.js