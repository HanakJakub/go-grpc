# GO GRPC solution

## Build and run without docker
Make sure you have install protoc and protoc go generator and run `make` command inside root directory. It will create all necessary files, build up server and client application and copy them to the root directory. After this you can start `./server` and after this `./client`. 
Server application will listen for incoming data from client, process data and send it back to the client. Client will do calculations, print the results and stop. Server will remain running and wait for another incoming data. 

## Build and run with docker

### Build docker image 
Make sure your pwd is root folder of repository and run build
```shell
docker build -t grpc-go-app .
```

### Run docker image and start server
Start docker image.
```shell
docker run -it -p 50005:50005 grpc-go-app
```

### Start client
Show list of containers where you should see grpc-go-app running.
```shell
docker ps
```

If docker is running you can start `./client` directly in root folder or enter container and run it there.

Copy grpc-go-app container id and replace [container-id]
```shell
docker exec -it [container-id] bash
```

This command will exec bash in the container. To start client and process the data start ./client
```shell
./client
```

## Testing
Tests are checking in Makefile, and for manual testing run `go test ./...` in root directory.