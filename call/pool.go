package call

import (
	"fmt"
	"midmsg/model"
	"sync"
	"time"
)


var TimeoutRequest sync.Pool
var AsyncReturn sync.Pool

func init()  {
	TimeoutRequest = sync.Pool{
		New: func() interface{} {
			b := model.CallInfo{}
			return &b
		},
	}

	AsyncReturn = sync.Pool{
		New: func() interface{} {
			b := model.AsyncReturnInfo{}
			return &b
		},
	}
}

func TestOtRequestPut(){

	t := model.CallInfo{
		Address:"192.168.0.1",
		Port:"5555",
		Timeout:time.Second*0,
		Service:"time.service",
		MsgBody:[]byte("time.service1"),
	}
	t1 := model.CallInfo{
		Address:"192.168.0.2",
		Port:"5555",
		Timeout:time.Second*0,
		Service:"time.service",
		MsgBody:[]byte("time.service2"),
	}
	t2 := model.CallInfo{
		Address:"192.168.0.3",
		Port:"5555",
		Timeout:time.Second*0,
		Service:"time.service",
		MsgBody:[]byte("time.service3"),
	}
	TimeoutRequest.Put(&t2)
	TimeoutRequest.Put(&t1)
	TimeoutRequest.Put(&t)


	TestOtRequestGet()
}

func TestOtRequestGet(){
	for{
		rq := TimeoutRequest.Get().(*model.CallInfo)
		if rq == nil {
			fmt.Println("rq == nil")
			return
		}else{
			//fmt.Println("rq.InBody:",rq.InBody==nil,"rq.Address:",rq.Address)
			if rq.MsgBody != nil && rq.Address != "" && rq.Port != "" {
				fmt.Println("取到了")
				fmt.Println("rq.InBody:",string(rq.MsgBody),"rq.Address:",rq.Address)
			}else{
				return
			}
		}
	}
}

func CallPoolRequest(){
	for{
		rq := TimeoutRequest.Get().(*model.CallInfo)
		if rq == nil {
			return
		}else{
			if rq.MsgBody != nil && rq.Address != "" && rq.Port != "" {
				CallClient(*rq,nil,nil)
			}else{
				return
			}
		}
	}
}

func CallPoolAsyncReturn(){
	for{
		rq := AsyncReturn.Get().(*model.AsyncReturnInfo)
		if rq == nil {
			return
		}else{
			if  rq.ClientIP != "" {
				AsyncReturnClient(rq)
			}else{
				return
			}
		}
	}
}

func PutPoolRequest(callinfo *model.CallInfo){
	TimeoutRequest.Put(callinfo)
}

func PutPoolAsyncReturn(returninfo *model.AsyncReturnInfo){
	AsyncReturn.Put(returninfo)
}

func TimerCallPool(){
	for  {
		select {
		case <- time.After(time.Second * 200):
			CallPoolRequest()
			CallPoolAsyncReturn()
		}
	}
}
