package cache

import (
	"google.golang.org/grpc"
	"midmsg/log"
	"sync"
)

var Concache sync.Map

func GetCache(key string)*grpc.ClientConn{
	conn,b := Concache.Load(key)
	if b {
		return conn.(*grpc.ClientConn)
	}else{
		return nil
	}
}


func StoreCache(key string, v *grpc.ClientConn){
	Concache.Store(key,v)
	log.Trace(Concache)
}