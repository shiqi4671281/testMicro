package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tencentyun/cos-go-sdk-v5"
	"image/png"
	"log"
	"nari/web/model"
	cappackage "nari/web/proto/captcha"
	userpackage "nari/web/proto/user"
	"nari/web/utils"
	"net/http"
	"net/url"
)

func init(){
	viper.AutomaticEnv()

	viper.SetConfigFile("./conf/config.yaml") // 指定配置文件
	//viper.AddConfigPath("./conf/")     // 指定查找配置文件的路径
	err := viper.ReadInConfig()        // 读取配置信息
	if err != nil {                    // 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// 监控配置文件变化
	viper.WatchConfig()
}

func GetSession(ctx *gin.Context){
	backcode:=make(map[string]interface{})
	s:=sessions.Default(ctx)
	userName:=s.Get("username")
	if userName!=nil{
		data:=make(map[string]interface{})
		data["name"]=userName
		backcode["errno"]=utils.RECODE_OK
		backcode["errmsg"]=utils.RecodeText(utils.RECODE_OK)
		backcode["data"]=data
	}else {
		backcode["errno"]=utils.RECODE_SESSIONERR
		backcode["errmsg"]=utils.RecodeText(utils.RECODE_SESSIONERR)
	}
	ctx.JSON(http.StatusOK,backcode)
}

func PostLogin(ctx *gin.Context){
	var logData struct{
		Mobile string `json:"mobile"`
		Password string `json:"password"`
	}
	_ = ctx.Bind(&logData)

	microClient:=utils.ConsulInit()
	client:=userpackage.NewUserService("user",microClient.Client())
	resp, err:= client.LoginUser(context.TODO(), &userpackage.LoginReq{Mobile: logData.Mobile, Password: logData.Password})
	if err!=nil{
		log.Fatal(err)
	}
	if resp.Errno==utils.RECODE_OK{
		s:=sessions.Default(ctx)
		s.Set("username",model.GetUsername(logData.Mobile))
		_ = s.Save()
	}
	ctx.JSON(http.StatusOK,resp)
}

func DeleteSession(ctx *gin.Context){
	backcode:=make(map[string]interface{})
	s:=sessions.Default(ctx)
	s.Delete("username")
	err:=s.Save()
	if err!=nil{
		backcode["errno"] = utils.RECODE_IOERR	// 没有合适错误,使用 IO 错误!
		backcode["errmsg"] = utils.RecodeText(utils.RECODE_IOERR)
	}else {
		backcode["errno"] = utils.RECODE_OK
		backcode["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	}
	ctx.JSON(http.StatusOK,backcode)
}

func GetImageCd(ctx *gin.Context){
	uuid:=ctx.Param("uuid")
	microClient:=utils.ConsulInit()

	client:=cappackage.NewCaptchaService("captcha",microClient.Client())

	rsp,err:=client.Call(context.TODO(),&cappackage.Request{Uuid:uuid})
	if err!=nil{
		log.Fatal(err)
	}
	var img captcha.Image
	_ = json.Unmarshal(rsp.Img, &img)

	_ = png.Encode(ctx.Writer, img)
}

func GetSmscd(ctx *gin.Context){
	phone :=ctx.Param("phone")
	imgcode:=ctx.Query("text")
	uuid:=ctx.Query("id")

	microClient:=utils.ConsulInit()
	client:=userpackage.NewUserService("user",microClient.Client())
	rsp,err:=client.SentSms(context.TODO(),&userpackage.Request{Phone:phone,Imgcode:imgcode,Uuid:uuid})
	if err!=nil{
		fmt.Println("调用user服务错误",err)
		return
	}
	ctx.JSON(http.StatusOK,rsp)
}

func PostRet(ctx *gin.Context){
	var regData struct{
		Mobile   string `json:"mobile"`
		Password string `json:"password"`
		SmsCode  string `json:"sms_code"`
	}
	_ = ctx.Bind(&regData)
	microClient:=utils.ConsulInit()
	client:=userpackage.NewUserService("user",microClient.Client())
	rsp, err := client.RegisterModel(context.TODO(),
		&userpackage.RegReq{Mobile: regData.Mobile,Password:regData.Password,Smscode:regData.SmsCode})
	if err!=nil{
		log.Fatal(err)
	}
	ctx.JSON(http.StatusOK,rsp)
}

func GetArea(ctx *gin.Context){

	microClient:=utils.ConsulInit()
	client:=userpackage.NewUserService("user",microClient.Client())
	rsp,err:=client.GetallArea(context.TODO(),&userpackage.Request{})

	if err!=nil{
		fmt.Println("user微服务查询失败")
		return
	}
	var areas []model.Area
	_ = json.Unmarshal(rsp.Data, &areas)
	resp:=make(map[string]interface{})
	resp["errno"]=rsp.Errno
	resp["errmsg"]=rsp.Errmsg
	resp["data"]=areas
	ctx.JSON(http.StatusOK,resp)
}

func GetUserInfo(ctx *gin.Context){
	microClient:=utils.ConsulInit()
	client:=userpackage.NewUserService("user",microClient.Client())
	backMsg:=make(map[string]interface{})
	s:=sessions.Default(ctx)
	userName:=s.Get("username")
	rsp,_:=client.GetUserInfo(context.TODO(),&userpackage.UserName{Username:userName.(string)})
	if userName!=nil{
		backMsg["errno"]=rsp.Errno
		backMsg["errmsg"]=rsp.Errmsg
		userData:=make(map[string]interface{})
		backData:=make(map[string]interface{})
		_ = json.Unmarshal(rsp.Data, &userData)
		backData["user_id"]=userData["ID"]
		backData["name"]=userData["Name"]
		backData["mobile"]=userData["Mobile"]
		backData["real_name"]=userData["Real_name"]
		backData["id_card"]=userData["Id_card"]
		backData["avatar_url"]=userData["Avatar_url"]
		backMsg["data"]=backData
		//fmt.Println(userData)
	}else {
		backMsg["errno"]=utils.RECODE_SESSIONERR
		backMsg["errmsg"]=utils.RecodeText(utils.RECODE_SESSIONERR)
	}
	ctx.JSON(http.StatusOK,backMsg)
}

func PutUserInfo(ctx *gin.Context){
	backMsg:=make(map[string]interface{})
	s:=sessions.Default(ctx)
	userName:=s.Get("username").(string)
	var newName struct{
		Name string `json:"name"`
	}
	_ = ctx.Bind(&newName)
	if len(userName)!=0{
		microClient:=utils.ConsulInit()
		client:=userpackage.NewUserService("user",microClient.Client())
		rsp,_:=client.UpdataUser(context.TODO(),&userpackage.UpdataName{Oldname:userName,Newname:newName.Name})
		backMsg["errno"]=rsp.Errno
		backMsg["errmsg"]=rsp.Errmsg
		backMsg["data"]=newName
		s.Set("username",newName.Name)
		_ = s.Save()

	}else {
		backMsg["errno"]=utils.RECODE_SESSIONERR
		backMsg["errmsg"]=utils.RecodeText(utils.RECODE_SESSIONERR)
	}
	ctx.JSON(http.StatusOK,backMsg)
}

func PostAvatar(ctx *gin.Context){
	backMsg:=make(map[string]interface{})
	file,_:=ctx.FormFile("avatar")
	cosurl:=viper.GetString("COSURL")
	//上传头像到腾讯云cos对象存储
	u, _ := url.Parse(cosurl)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  viper.GetString("SECRETID"),
			SecretKey: viper.GetString("SECRETKEY"),
		},
	})
	// 对象键（Key）是对象在存储桶中的唯一标识。
	name := "avatar/"+file.Filename
	urlCom:=cosurl+"/"+name
	fd,err:=file.Open()
	if err != nil {
		panic(err)
	}
	_, err = c.Object.Put(context.Background(), name, fd, nil)
	if err != nil {
		panic(err)
	}

	userName:=sessions.Default(ctx).Get("username")
	er:=model.UpdataAvatar(userName.(string),urlCom)
	if er!=nil{
		backMsg["errno"]=utils.RECODE_DBERR
		backMsg["errmsg"]=utils.RecodeText(utils.RECODE_DATAERR)
	}else {
		backMsg["errno"]=utils.RECODE_OK
		backMsg["errmsg"]=utils.RecodeText(utils.RECODE_OK)
		backMsg["data"]=cosurl+"/"+name
	}
	ctx.JSON(http.StatusOK,backMsg)
}
