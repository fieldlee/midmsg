package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"midmsg/call"
	"midmsg/handle"
	"midmsg/log"
	pb "midmsg/proto"
	"midmsg/utils"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
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
	go func() {
		log.Info("go routine waiting shutdown...")
		waitForShutdown(rpcServer)
	}()
	pb.RegisterMidServiceServer(rpcServer, &handle.MsgHandle{})
	reflection.Register(rpcServer)
	if err = rpcServer.Serve(listener); err != nil {
		log.Fatal("failed serve at: " + Host + ":" + Port)
	}
}

func waitForShutdown(srv *grpc.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// Block until we receive our signal.
	<-interruptChan
	// graceful shutdown grpc service
	srv.GracefulStop()

	log.Info("grpc stop...")
	// shutdown service
	os.Exit(0)
}