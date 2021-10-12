package main

import (
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/registry/consul"
	"github.com/spf13/viper"
	"captcha/handler"
	cappackage "captcha/proto/captcha"
)

func main() {

	viper.AutomaticEnv()
	viper.SetConfigFile("./conf/config.yaml") // 指定配置文件
	//viper.AddConfigPath("./conf/")     // 指定查找配置文件的路径
	err := viper.ReadInConfig()        // 读取配置信息
	if err != nil {                    // 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// 监控配置文件变化
	viper.WatchConfig()

	consulReg:=consul.NewRegistry(func(options *registry.Options) {
		options.Addrs=[]string{
			viper.GetString("CONSUL"),
		}
	})
	// New Service
	service := micro.NewService(
		micro.Address(viper.GetString("HOSTPORT")),
		micro.Name("captcha"),
		micro.Registry(consulReg),
		micro.Version("latest"),
	)

	// Initialise service
	//service.Init()

	// Register Handler
	_ = cappackage.RegisterCaptchaHandler(service.Server(), new(handler.Captcha))


	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
