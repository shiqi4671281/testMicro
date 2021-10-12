package main

import (
	_"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Student struct {
	gorm.Model
	Name string
	Score int
}

var GlobalConn *gorm.DB

func main(){
	conn,err:=gorm.Open("mysql","root:hello@123@tcp(127.0.0.1:3306)/shiqi?parseTime=True&loc=Local")
	if err!=nil{
		log.Fatal(err)
	}
	//定义全局连接池
	GlobalConn=conn
	//不使用表复数
	GlobalConn.SingularTable(true)
	//创建表
	//GlobalConn.AutoMigrate(new(Student))
	//插入数据
	//var stu2 Student
	//GlobalConn.Select("name,score").Where("id=?",4).Find(&stu2)
	//fmt.Println(stu2)
	//GlobalConn.DropTableIfExists(&stu)

	//GlobalConn.Model(new(Student)).Where("id=?",3).Update("deleted_at","NULL")

	GlobalConn.Unscoped().Where("name=?", "zhang3").Delete(new(Student))
	//GlobalConn.Unscoped().Delete(new(Student)).Where("name=?","shiqi")
}
