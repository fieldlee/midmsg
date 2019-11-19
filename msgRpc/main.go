package main

import (
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"midmsg/utils"
	"google.golang.org/grpc"
	pb "midmsg/proto"
)

var (
	Host = utils.Address
	Port = string(utils.Port)
)

func main()  {
	listener, err := net.Listen("tcp", Host+":"+Port)
	if err != nil {
		log.Fatalln("faile listen at: " + Host + ":" + Port)
	} else {
		log.Println("Demo server is listening at: " + Host + ":" + Port)
	}
	rpcServer := grpc.NewServer()
	pb.RegisterMidServiceServer(rpcServer, &FormatData{})
	reflection.Register(rpcServer)
	if err = rpcServer.Serve(listener); err != nil {
		log.Fatalln("faile serve at: " + Host + ":" + Port)
	}
}