package main

import (
	"fmt"
	pb "github.com/darshankapadiya19/grpc-gw-rest-proto/gen"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
)

func main() {
	// Create a HelloRequest message
	req := &pb.HelloRequest{Name: "Darshan"}

	// Serialize the message to binary format
	data, _ := proto.Marshal(req)

	// Write binary data to proto.bin file
	file, _ := os.Create("proto.bin")
	defer file.Close()

	file.Write(data)

	log.Println("Data written to proto.bin")

	// Read from file and unmarshal
	file, _ = os.Open("proto.bin")
	defer file.Close()

	fileData, _ := os.ReadFile("proto.bin")

	x := &pb.HelloRequest{}
	if err := proto.Unmarshal(fileData, x); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshalling request: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Data name: %s", x.Name)
}
