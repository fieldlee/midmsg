package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"midmsg/handle"
	pb "midmsg/proto"
	"midmsg/utils"
	"net"
	"os"
)

var (
	Host = utils.Address
	Port = fmt.Sprintf("%d",utils.Port)
)

func main()  {
	fmt.Println(Port)
	listener, err := net.Listen("tcp", Host+":"+Port)
	if err != nil {
		log.Fatalln("faile listen at: " + Host + ":" + Port)
	} else {
		log.Println("server is listening at: " + Host + ":" + Port)
	}
	rpcServer := grpc.NewServer()
	pb.RegisterMidServiceServer(rpcServer, &handle.MsgHandle{})
	reflection.Register(rpcServer)
	//test()
	if err = rpcServer.Serve(listener); err != nil {
		log.Fatalln("faile serve at: " + Host + ":" + Port)
	}
}

func test(){
	//启动多线程处理
	fmt.Println("=================")
	d := handle.NewDispatcher(utils.MaxWorker)
	d.Run()
	fmt.Println("=================")
	t := &handle.MsgHandle{}
	body := getbody()
	for i:=0;i<1000000 ; i++ {
		fmt.Println(i)
		tbody := &pb.NetReqInfo{
			M_Body:body,
		}
		rsp,err := t.Sync(context.TODO(),tbody)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(rsp)
	}
}
func getbody()[]byte{
	bodyByte := []byte{}
	fileName := "./1.txt"
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return bodyByte
	}
	defer file.Close()

	buf := bufio.NewReader(file)

	var i = 0
	for {
		line, err := buf.ReadBytes('\n')
		if i == 1 {
			line = line[:]
			bodyByte = line
			return bodyByte
		}

		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return bodyByte
			}
		}
		i ++
	}
	return bodyByte
}