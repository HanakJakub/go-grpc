
all: test client server

dep:
	@echo "Install dependencies"
	cd src && glide install

protoc:
	@echo "Generating Go files"
	cd src/protobuf && protoc --go_out=plugins=grpc:. *.proto

test:
	go test ./...

server: protoc
	@echo "Building server"
	go build -o server \
		./src/server
		
client: protoc
	@echo "Building client"
	go build -o client \
		./src/client

clean:
	go clean ./...
	rm -f server client

.PHONY: client server protoc dep test
