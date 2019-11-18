package main

import (

	"log"
	"net"

	"google.golang.org/grpc"

)

const (
	PORT = "9002"
)

func main() {
	server := grpc.NewServer()

	pb.RegisterStreamServiceServer(server, &StreamService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	server.Serve(lis)
}


