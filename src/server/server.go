package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"go-grpc/src/lib"
	pb "go-grpc/src/protobuf"

	"google.golang.org/grpc"
)

var (
	filepath, filename string
)

type server struct{}

func (s server) Purchaser(srv pb.Processor_PurchaserServer) error {

	log.Println("start new server")
	ctx := srv.Context()

	for {
		// exit if context is done
		// or continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// receive data from stream
		req, err := srv.Recv()
		if err == io.EOF {
			// return will close stream from server side
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}

		// skip not existing files
		if _, err := os.Stat(req.GetPath()); os.IsNotExist(err) {
			log.Println("File " + req.GetPath() + " does not exist")
			continue
		}

		// load user file and parse data to struct
		f := loadUserFile(req.GetPath() + req.GetName())
		user := parseUser(req.GetName(), f)

		unique := make(map[string]int64)
		for _, p := range user.Purchases {
			if p.ShouldProcess() && isItUniqueForUser(unique, p) {
				// convert purchase to string
				pString, err := p.ToString()
				if err != nil {
					log.Printf("receive error %v", err)
					continue
				}

				// prepare response to client
				resp := pb.Response{Result: pString}

				// send response
				if err := srv.Send(&resp); err != nil {
					log.Printf("send error %v", err)
				}

				// print sent data and file from which data came
				log.Println(req.GetName(), p)
			}
		}
	}
}

func main() {
	// create listener for port 50005
	lis, err := net.Listen("tcp", ":50005")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create grpc server
	s := grpc.NewServer()
	pb.RegisterProcessorServer(s, server{})

	// start server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// loadUserFile from a given path
func loadUserFile(path string) []byte {

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

// parseUser from a file to a user struct
func parseUser(name string, file []byte) (u lib.User) {
	u.ID = getUserIDFromFileName(name)

	err := json.Unmarshal(file, &u)
	if err != nil {
		log.Fatal(err)
	}

	return
}

// getUserIDFromFileName will get user ID from given file
func getUserIDFromFileName(name string) (id uint64) {
	idx := strings.Index(name, ".")

	id, err := strconv.ParseUint(name[:idx], 10, 64)
	if err != nil {
		log.Println(err)
	}

	return
}

// isItUniqueForUser returns true if it is first occurance of given type of purchase
func isItUniqueForUser(unique map[string]int64, p lib.Purchase) bool {
	if _, ok := unique[p.Type]; ok {
		return false
	}

	unique[p.Type] = p.Amount

	return true
}
