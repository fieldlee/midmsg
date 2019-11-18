package main

import (

	"log"
	"net"
	context "context"
	"google.golang.org/grpc"

	pb "midmsg/proto"
)

type MidService struct{}

const (
	PORT = "9002"
)

func main() {
	server := grpc.NewServer()

	pb.RegisterMidServiceServer(server,&MidService{})

	listener, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	server.Serve(listener)
}

func (mid *MidService)Sync(ctx context.Context,req *pb.GJ_Net_Pack)(*pb.GJ_Return_Pack,error){
	return nil,nil
}

func (mid *MidService)Async(ctx context.Context,req *pb.GJ_Net_Pack)(*pb.GJ_Return_Pack,error){
	return nil,nil
}