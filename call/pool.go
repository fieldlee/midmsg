package call

import (
	"midmsg/model"
	"sync"
	"time"
)

var TimeoutRequest CallInfoPool
var AsyncReturn    AsyncReturnPool
var AsyncAnswer	   sync.Map

type CallInfoPool struct {
	mux sync.RWMutex
	CallInfoList []model.CallInfo
}

type AsyncReturnPool struct {
	mux sync.RWMutex
	AsyncReturnPool []model.CallInfo
}

func init()  {
	TimeoutRequest = CallInfoPool{
		CallInfoList:make([]model.CallInfo,0),
	}
	AsyncReturn = AsyncReturnPool{
		AsyncReturnPool:make([]model.CallInfo,0),
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
	p.AsyncReturnPool = make([]model.CallInfo,0)
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

func (p *AsyncReturnPool)PutPoolAsyncReturn(returninfo model.CallInfo){
	p.mux.Lock()
	list := p.AsyncReturnPool
	p.AsyncReturnPool = append(list,returninfo)
	p.mux.Unlock()
}

func CheckAsyncAnswer(key string) bool{
	_,ok := AsyncAnswer.Load(key)
	return ok
}

func StoreAsyncAnswer(key string,v interface{}){
	AsyncAnswer.Store(key,v)
}

func LoadAsyncAnswer(key string)model.CallInfo{
	v,ok := AsyncAnswer.Load(key)
	if ok {
		defer AsyncAnswer.Delete(key)
		return v.(model.CallInfo)
	}
	return model.CallInfo{}
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
