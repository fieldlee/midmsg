package call

import (
	"midmsg/model"
	"sync"
	"time"
)

var TimeoutRequest CallInfoPool
var AsyncReturn    AsyncReturnPool

type CallInfoPool struct {
	mux sync.RWMutex
	CallInfoList []model.CallInfo
}

type AsyncReturnPool struct {
	mux sync.RWMutex
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


func (c *CallInfoPool)CallPoolRequest(){

	c.mux.Lock()
	list := c.CallInfoList
	c.CallInfoList = make([]model.CallInfo,0)
	c.mux.Unlock()

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

func (c *CallInfoPool)PutPoolRequest(callinfo model.CallInfo){
	c.mux.Lock()
	list := c.CallInfoList
	c.CallInfoList = append(list,callinfo)
	c.mux.Unlock()
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
		case <- time.After(time.Second * 200):
			TimeoutRequest.CallPoolRequest()
			AsyncReturn.CallPoolAsyncReturn()
		}
	}
}
