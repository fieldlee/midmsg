package utils

import (
	"log"
	"github.com/spf13/viper"
	"reflect"
)
var (
	Con *Config
	Address string
	Port uint32

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
	Con.V.AddConfigPath("../")
	//设置配置文件类型
	Con.V.SetConfigType("yaml")
	if err := Con.V.ReadInConfig(); err != nil {
		log.Fatal(err.Error())
	}
	return Con
}

func init()  {
	Con = InitConfig()
}

func GetMidAddress()string{
	return Con.V.GetString("address")
}

func GetMidPort()int32{
	return Con.V.GetInt32("port")
}

func GetServices(key string)map[string]string{
	//serviceMap := make(map[string]string)
	services := Con.V.GetStringMap("services")

	for k , v := range services{
		if k == key {
			if reflect.TypeOf(v).Kind() == reflect.Map {
				return v.(map[string]string)
			}
		}
	}
	return nil
}