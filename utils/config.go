package utils

import (
	"github.com/spf13/viper"
	"midmsg/log"
	"reflect"
)

var (
	Con *Config
	Address string
	Port uint32
	ClientPort uint32
	Services map[string]interface{}
	SubScribe map[string]interface{}
	MaxQueue uint32
	MaxWorker uint32
)
type Config struct {
	V *viper.Viper
}
func InitConfig () *Config {
	Con := &Config{
		V: viper.New(),
	}
	//设置配置文件的名字
	Con.V.SetConfigName("config")
	//添加配置文件所在的路径,注意在Linux环境下%GOPATH要替换为$GOPATH
	Con.V.AddConfigPath("./")
	//设置配置文件类型
	Con.V.SetConfigType("yaml")
	if err := Con.V.ReadInConfig(); err != nil {
		log.Fatal(err.Error())
	}
	return Con
}

func init()  {
	ReloadConfig()
}

func ReloadConfig(){
	Con = InitConfig()
	Address = GetMidAddress()
	Port = GetMidPort()
	ClientPort = GetClientPort()
	Services = GetServices()
	MaxWorker = GetMaxWorker()
	MaxQueue = GetMaxQueue()
	SubScribe = GetSubscribe()
}

func GetMidAddress()string{
	return Con.V.GetString("address")
}

func GetMidPort()uint32{
	return Con.V.GetUint32("port")
}

func GetClientPort()uint32{
	return Con.V.GetUint32("clientport")
}

func GetMaxWorker()uint32{
	return Con.V.GetUint32("maxwoker")
}

func GetMaxQueue()uint32{
	return Con.V.GetUint32("maxqueue")
}

func GetServices()map[string]interface{}{
	//serviceMap := make(map[string]string)
	services := Con.V.GetStringMap("services")
	return services
}

func GetSubscribe()map[string]interface{}{
	//serviceMap := make(map[string]string)
	subscribe := Con.V.GetStringMap("subscribe")
	return subscribe
}

func GetServiceByKey(key string)map[string]interface{}{
	for k , v := range Services{
		if k == key {
			if reflect.TypeOf(v).Kind() == reflect.Map {
				return v.(map[string]interface{})
			}
		}
	}
	return nil
}

func GetSubscribeByKey(key string)[]interface{}{
	for k,v := range SubScribe{
		if k == key {
			log.Info(v)
			if reflect.TypeOf(v).Kind() == reflect.Map {
				SubMap := v.(map[string]interface{})
				addrs := SubMap["subaddrs"]
				if reflect.TypeOf(addrs).Kind() == reflect.Slice {
					return addrs.([]interface{})
				}
			}
		}
	}
	return nil
}