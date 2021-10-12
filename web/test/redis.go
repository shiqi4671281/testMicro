package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
)

func main() {
	conn,err:=redis.Dial("tcp","127.0.0.1:6379")
	if err!=nil{
		log.Fatal(err)
	}
	defer  conn.Close()

	code,err:=redis.String(conn.Do("get","shiqi"))
	if err!=nil{
		fmt.Println("验证码错误")
		return
	}else {
		fmt.Println(code)
	}
}

