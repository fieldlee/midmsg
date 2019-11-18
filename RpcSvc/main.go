package main

import (

	"log"
	"net"
	"context"
	"google.golang.org/grpc"
	pb "midmsg/proto"
	"midmsg/RpcSvc/service"
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
	reps :=  new(pb.GJ_Return_Pack)
	err := service.SyncService(ctx,req)
	if err != nil {
		return reps,err
	}
	return reps, nil
}

func (mid *MidService)Async(ctx context.Context,req *pb.GJ_Net_Pack)(*pb.GJ_Return_Pack,error){
	reps :=  new(pb.GJ_Return_Pack)
	err := service.AsyncService(ctx,req)
	if err != nil {
		return reps,err
	}
	return reps, nil
}