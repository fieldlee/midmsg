package service

import (
	"context"
	"log"
	pb "midmsg/proto"
	"github.com/spf13/viper"
)

type Config struct {
	V *viper.Viper
}


func InitConfig () *Config {
	Con := &Config{
		V:viper.New(),
	}
	//设置配置文件的名字
	Con.V.SetConfigName("config")
	//添加配置文件所在的路径,注意在Linux环境下%GOPATH要替换为$GOPATH

	Con.V.AddConfigPath("../../")
	//设置配置文件类型
	Con.V.SetConfigType("yaml")
	if err := Con.V.ReadInConfig(); err != nil{
		log.Fatal(err.Error())
	}
	return Con
}

var config *Config

func init()  {
	config = InitConfig()
}

func SyncService(ctx context.Context,req *pb.GJ_Net_Pack)error{
	mapReq := req.GetM_Net_Pack()
	///// 多个客户
	for k,v := range mapReq{
		go RequestByClient(k,v.M_Msg,v.M_MsgBody)
	}

	return nil
}


func RequestByClient(identity uint32 ,msg []byte,msgbody *pb.Min_Net_MsgBody){
	//ip := config.V.GetString("")
	//port := config.V.GetInt("")


}


func AsyncService(ctx context.Context,req *pb.GJ_Net_Pack)error{

	return nil
}