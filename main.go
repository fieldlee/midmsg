package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"midmsg/call"
	"midmsg/handle"
	"midmsg/log"
	"runtime"

	//"midmsg/model"
	pb "midmsg/proto"
	"midmsg/utils"
	"net"
	"os"
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

func test(){
	//启动多线程处理
	body := getbody()
	t := &handle.MsgHandle{}

	go func() {
		for i:=0; i < 10000 ; i++ {
			fmt.Println(i)
			tbody := &pb.NetReqInfo{
				M_Body:body,
			}
			rsp,err := t.Sync(context.TODO(),tbody)
			if err != nil {
				log.Error(err.Error())
			}
			log.Info(rsp)
			//handleBody := handle.HandleBody{
			//	M_Body:body,
			//}
			//handle.JobQueue <- handleBody
		}

		//for {
		//	select {
		//	case _ = <-JobDone:
		//		fmt.Println("Done Job")
		//	}
		//}
	}()


}
func getbody()[]byte{
	fileName := "./2.txt"
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		log.Error(err.Error())
		return nil
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	bodyByte := make([]byte,110)
	_,err = buf.Read(bodyByte)
	if err != nil {
		return nil
	}
	return bodyByte
}