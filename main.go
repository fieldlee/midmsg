package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"midmsg/call"
	"midmsg/handle"
	"midmsg/log"
	"runtime"

	pb "midmsg/proto"
	"midmsg/utils"
	"net"
)

var (
	Host = utils.Address
	Port = fmt.Sprintf("%d",utils.Port)
)

func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main()  {

	log.SetLogLevel(logrus.TraceLevel)
	log.SetLogFormatter(&logrus.TextFormatter{DisableColors:true})

	d := handle.NewDispatcher(utils.MaxWorker,handle.JobDone)
	d.Run()
	////////启动监听sync pool里的数据，并200秒发送一次
	go call.TimerCallPool()

	listener, err := net.Listen("tcp", Host+":"+Port)
	if err != nil {
		log.Fatal("failed listen at: " + Host + ":" + Port)
	} else {
		log.Info("server is listening at: " + Host + ":" + Port)
	}
	rpcServer := grpc.NewServer()
	pb.RegisterMidServiceServer(rpcServer, &handle.MsgHandle{})
	reflection.Register(rpcServer)
	if err = rpcServer.Serve(listener); err != nil {
		log.Fatal("failed serve at: " + Host + ":" + Port)
	}
}
