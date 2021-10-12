package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"log"
)

var RedisPool redis.Pool

func InitRedis(){
	RedisPool=redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp",viper.GetString("REDIS"))
		},
		MaxIdle:         2,
		MaxActive:       3,
		IdleTimeout:     60,
		MaxConnLifetime: 60*5,
	}
}

func CheckCapt(imgcode,uuid string )  bool{
	conn:=RedisPool.Get()
	defer  conn.Close()

	code,err:=redis.String(conn.Do("get",uuid))
	if err!=nil{
		log.Println(err)
		return false
	}
	return  code==imgcode
}

func InputSms(phone,smscode string) {
	conn:=RedisPool.Get()
	defer conn.Close()
	_, _ = conn.Do("setex", phone+"_code", 60*3, smscode)
}

func CheckSms(mobile,smscode string) bool{
	conn:=RedisPool.Get()
	defer conn.Close()

	code,err:=redis.String(conn.Do("get",mobile+"_code"))
	if err!=nil{
		fmt.Println(err)
	}
	return code==smscode
}

func GetUserInfo(username string) (User,error){
	var user User
	err:=Gdb.Where("name=?",username).Find(&user).Error
	return user,err
}

func UpdataUser(oldName, newName string) error {
	err:=Gdb.Model(new(User)).Where("name=?",oldName).Update("name",newName).Error
	return err
}
