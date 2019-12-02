package call

import (
	"github.com/flyaways/pool"
	"google.golang.org/grpc"
	"midmsg/log"
	"sync"
	"time"
)


var Concache sync.Map

func GetCache(key string)*pool.GRPCPool{
	conn,b := Concache.Load(key)
	if b {
		return conn.(*pool.GRPCPool)
	}else{
		return nil
	}
}

func StoreCache(key string, v *pool.GRPCPool)*pool.GRPCPool{
	gPool,loaded := Concache.LoadOrStore(key,v)
	if loaded {
		return gPool.(*pool.GRPCPool)
	}else{
		return nil
	}
}


func NewPool(addr string)*pool.GRPCPool{
	options := &pool.Options{
		InitTargets:  []string{addr},
		InitCap:      5,
		MaxCap:       30,
		DialTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 60,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}
	//初始化连接池
	p, err := pool.NewGRPCPool(options, grpc.WithInsecure())

	if err != nil {
		log.Error(err)
		return nil
	}
	if p == nil {
		return nil
	}

	g := StoreCache(addr,p)
	if g == nil {
		return nil
	}

	return g
}