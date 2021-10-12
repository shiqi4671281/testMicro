package model

import (
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"log"
)

func InputCapt(uuid,code string ){
	conn, err:= redis.Dial("tcp", viper.GetString("REDIS"))
	if err!=nil{
		log.Fatal(err)
	}
	defer conn.Close()

	_, _ = conn.Do("setex", uuid, 60*3, code)
}

