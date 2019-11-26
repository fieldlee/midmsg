package call

import (
	"midmsg/model"
	"sync"
	"time"
)

var count = 0
var TimeoutRequest CallInfoPool
var AsyncReturn    AsyncReturnPool

type CallInfoPool struct {
	mux sync.Mutex
	CallInfoList []model.CallInfo
}

type AsyncReturnPool struct {
	mux sync.Mutex
	AsyncReturnPool []model.AsyncReturnInfo
}

func init()  {
	TimeoutRequest = CallInfoPool{
		CallInfoList:make([]model.CallInfo,0),
	}
	AsyncReturn = AsyncReturnPool{
		AsyncReturnPool:make([]model.AsyncReturnInfo,0),
	}
}


func CallPoolRequest(){
	TimeoutRequest.mux.Unlock()

	list := TimeoutRequest.CallInfoList
	TimeoutRequest.CallInfoList = make([]model.CallInfo,0)
	TimeoutRequest.mux.Lock()

	for _,v := range list {
		CallClient(v,nil,nil)
	}
}

func (p *AsyncReturnPool) CallPoolAsyncReturn(){
	p.mux.Lock()
	list := p.AsyncReturnPool
	p.AsyncReturnPool = make([]model.AsyncReturnInfo,0)
	p.mux.Unlock()

	for _,v := range list {
		AsyncReturnClient(v)
	}
}

func PutPoolRequest(callinfo model.CallInfo){
	TimeoutRequest.mux.Unlock()
	list := TimeoutRequest.CallInfoList
	TimeoutRequest.CallInfoList = append(list,callinfo)
	TimeoutRequest.mux.Lock()
}

func (p *AsyncReturnPool)PutPoolAsyncReturn(returninfo model.AsyncReturnInfo){
	p.mux.Lock()
	list := p.AsyncReturnPool
	p.AsyncReturnPool = append(list,returninfo)
	p.mux.Unlock()
}

func TimerCallPool(){
	for  {
		select {
		case <- time.After(time.Second * 20):
			CallPoolRequest()
			AsyncReturn.CallPoolAsyncReturn()
		}
	}
}
