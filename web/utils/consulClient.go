package utils

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/consul"
	"github.com/spf13/viper"
)

func ConsulInit() micro.Service{
	consulReg:=consul.NewRegistry(func(options *registry.Options) {
		options.Addrs=[]string{
			viper.GetString("CONSUL"),
		}
	})

	microClient:=micro.NewService(micro.Registry(consulReg))
	return microClient
}
