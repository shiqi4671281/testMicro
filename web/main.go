package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"web/controller"
	"web/model"
)

func LoginSession() gin.HandlerFunc{
	return func (ctx *gin.Context){
		s:=sessions.Default(ctx)
		userName:=s.Get("username")
		if userName==nil{
			ctx.Abort()
		}else {
			ctx.Next()
		}

	}
}



func main(){
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



	router:=gin.Default()
	store,_:=redis.NewStore(2,"tcp",viper.GetString("REDIS"),viper.GetString("PASSWORD"),[]byte("shiqi"))
	router.Use(sessions.Sessions("loginsession",store))
	router.Static("/home","view")
	router.GET("/", func(ctx *gin.Context) {
		ctx.Request.URL.Path="/home"
		router.HandleContext(ctx)
	})

	r1:=router.Group("api/v1.0")
	{
		r1.POST("sessions",controller.PostLogin)
		r1.GET("/imagecode/:uuid",controller.GetImageCd)
		r1.GET("/smscode/:phone",controller.GetSmscd)
		r1.POST("/users",controller.PostRet)
		r1.GET("/areas",controller.GetArea)

		LoginSession()
		r1.GET("/user",controller.GetUserInfo)
		r1.PUT("/user/name",controller.PutUserInfo)
		r1.GET("/session",controller.GetSession)
	    r1.DELETE("/session",controller.DeleteSession)
		r1.POST("user/avatar",controller.PostAvatar)
	}


	_ = router.Run(viper.GetString("HOSTPORT"))
}
