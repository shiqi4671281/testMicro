package main

import (
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/registry/consul"
	"github.com/spf13/viper"
	"user/handler"
	"user/model"
	userpackage "user/proto/user"
)



func main() {
	model.InitRedis()
	_, err := model.InitDb()
	if err!=nil{
		log.Fatal(err)
	}


	viper.AutomaticEnv()
	viper.SetConfigFile("./conf/config.yaml") // 指定配置文件
	//viper.AddConfigPath("./conf/")     // 指定查找配置文件的路径
	err = viper.ReadInConfig()        // 读取配置信息
	if err != nil {                    // 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// 监控配置文件变化
	viper.WatchConfig()

	// New Service
	consulReg:=consul.NewRegistry(func(options *registry.Options) {
		options.Addrs=[]string{
			viper.GetString("CONSUL"),
		}
	})
	service := micro.NewService(
		micro.Address(viper.GetString("HOSTPORT")),
		micro.Registry(consulReg),
		micro.Name("user"),
		micro.Version("latest"),
	)

	// Register Handler
	_ = userpackage.RegisterUserHandler(service.Server(), new(handler.User))


	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
