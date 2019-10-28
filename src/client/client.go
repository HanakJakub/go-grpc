package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"go-grpc/src/lib"
	pb "go-grpc/src/protobuf"

	"google.golang.org/grpc"
)

// set the directory for json data
const dataPath = "./task/data/"

var (
	wg  sync.WaitGroup
	sum int64
)

type purchasesSorter []lib.Purchase

func (a purchasesSorter) Len() int           { return len(a) }
func (a purchasesSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a purchasesSorter) Less(i, j int) bool { return a[i].Amount < a[j].Amount }

func main() {
	start := time.Now()

	c := retrieveDataFromListener()
	counter := groupByType(c)

	res := anonymizeData(counter)

	log.Println("Median is: ", median(res))
	log.Println("Avegare is: ", mean(res))

	elapsed := time.Since(start)
	log.Printf("Process took %s", elapsed)
}

// retrieveDataFromListener will get all data from server
func retrieveDataFromListener() chan *pb.Response {
	results := make(chan *pb.Response)

	go func() {
		wg.Add(1)

		// dial server
		conn, err := grpc.Dial(":50005", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("can not connect with server %v", err)
		}

		// create new stream
		client := pb.NewProcessorClient(conn)
		stream, err := client.Purchaser(context.Background())
		if err != nil {
			log.Fatalf("openn stream error %v", err)
		}

		// run walk in goroutine and send requests to server
		go func() {
			filepath.Walk(dataPath, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					name := info.Name()
					if filepath.Ext(name) == ".json" {
						req := pb.Request{Path: dataPath, Name: name}
						if err := stream.Send(&req); err != nil {
							log.Fatalf("can not send %v", err)
						}
					}
				}

				return err
			})

			if err := stream.CloseSend(); err != nil {
				log.Println(err)
			}

		}()

		// goroutine to listen for any data that was send
		// and populate results channel
		go func() {
			for {
				resp, err := stream.Recv()
				if err == io.EOF {
					wg.Done()
					return
				}
				if err != nil {
					log.Fatalf("can not receive %v", err)
				}

				results <- resp
			}
		}()

		// run checker and close channel
		go func() {
			wg.Wait()
			close(results)
		}()
	}()

	return results
}

// median will calculate median of all purchases Amounts
func median(p purchasesSorter) float64 {
	sort.Sort(purchasesSorter(p))

	if len(p)%2 == 0 {
		middle := len(p) / 2
		higher := float64(p[middle].Amount)
		lower := float64(p[middle-1].Amount)
		return (higher + lower) / 2
	}

	middle := len(p) / 2
	return float64(p[middle].Amount)
}

// mean will calculate mean of all purchases Amounts
func mean(p purchasesSorter) float64 {
	var sum float64

	for _, row := range p {
		sum += float64(row.Amount)
	}

	return sum / float64(len(p))
}

// groupByType will group all responses by type
func groupByType(c chan *pb.Response) map[string][]*pb.Response {
	counter := make(map[string][]*pb.Response, len(c))

	for i := range c {
		p := lib.Purchase{}

		err := json.Unmarshal([]byte(i.GetResult()), &p)
		if err != nil {
			log.Fatal(err)
		}

		counter[p.Type] = append(counter[p.Type], i)
	}

	return counter
}

// anonymizeData will make anonymization of data and return purchaseSorter
func anonymizeData(counter map[string][]*pb.Response) (res purchasesSorter) {
	for _, row := range counter {
		if len(row) > 5 {
			for _, j := range row {
				p := lib.Purchase{}

				err := json.Unmarshal([]byte(j.GetResult()), &p)
				if err != nil {
					log.Fatal(err)
				}

				res = append(res, p)
			}
		}
	}

	return res
}
