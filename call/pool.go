package call

import (
	"fmt"
	"sync"
	"time"
)
var TimeoutRequest sync.Pool

type OtRequest struct{
	Address string
	Port 	string
	Timeout time.Duration
	Service string
	InBody  []byte
}

func init()  {
	TimeoutRequest = sync.Pool{
		New: func() interface{} {
			b := OtRequest{}
			return &b
		},
	}
}

func TestOtRequestPut(){

	t := OtRequest{
		Address:"192.168.0.1",
		Port:"5555",
		Timeout:time.Second*0,
		Service:"time.service",
		InBody:[]byte("time.service1"),
	}
	t1 := OtRequest{
		Address:"192.168.0.2",
		Port:"5555",
		Timeout:time.Second*0,
		Service:"time.service",
		InBody:[]byte("time.service2"),
	}
	t2 := OtRequest{
		Address:"192.168.0.3",
		Port:"5555",
		Timeout:time.Second*0,
		Service:"time.service",
		InBody:[]byte("time.service3"),
	}
	TimeoutRequest.Put(&t2)
	TimeoutRequest.Put(&t1)
	TimeoutRequest.Put(&t)


	TestOtRequestGet()
}

func TestOtRequestGet(){
	for{
		rq := TimeoutRequest.Get().(*OtRequest)
		if rq == nil {
			fmt.Println("rq == nil")
			return
		}else{
			//fmt.Println("rq.InBody:",rq.InBody==nil,"rq.Address:",rq.Address)
			if rq.InBody != nil && rq.Address != "" && rq.Port != "" {
				fmt.Println("取到了")
				fmt.Println("rq.InBody:",string(rq.InBody),"rq.Address:",rq.Address)
			}else{
				return
			}
		}
	}
}

func CallPoolRequest(){
	for{
		rq := TimeoutRequest.Get().(*OtRequest)
		if rq == nil {
			return
		}else{
			if rq.InBody != nil && rq.Address != "" && rq.Port != "" {
				CallClient(rq.Address,rq.Port,rq.Timeout,rq.Service,rq.InBody)
			}else{
				return
			}
		}
	}
}

func PutPoolRequest(address,port string,timeout time.Duration,service string,msg []byte){
	t := OtRequest{
		Address:address,
		Port:port,
		Timeout:timeout,
		Service:service,
		InBody:msg,
	}
	TimeoutRequest.Put(&t)
}

func TimerCallPool(){
	for  {
		select {
		case <- time.After(time.Second * 2):
			CallPoolRequest()
		}
	}
}
