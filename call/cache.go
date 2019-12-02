package call

import (
	"github.com/flyaways/pool"
	"time"
	"sync"
	"midmsg/log"
	"google.golang.org/grpc"
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

func StoreCache(key string, v *pool.GRPCPool){
	Concache.Store(key,v)
	log.Trace(Concache)
}


func New(addr string)*pool.GRPCPool{
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
	return p
}