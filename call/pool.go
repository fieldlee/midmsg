package call

import (
	"fmt"
	"time"
)



func CallPoolRequest(){
	fmt.Println("call pool request")
	for{
		rq := TimeoutRequest.Get().(*ToutRequest)
		if rq == nil {
			fmt.Println("rq == nil")
			return
		}else{
			fmt.Println("rq.InBody:",rq.InBody==nil,"rq.Address:",rq.Address)
			if rq.InBody != nil && rq.Address != "" && rq.Port != "" {
				CallClient(rq.Address,rq.Port,&rq.Timeout,rq.Service,rq.InBody)
			}else{
				return
			}
		}
	}
}

func PutPoolRequest(address,port string,timeout *time.Duration,service string,msg []byte){
	t := ToutRequest{
		Address:address,
		Port:port,
		Timeout:*timeout,
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
