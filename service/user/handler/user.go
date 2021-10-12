package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"math/rand"
	"nari/service/user/model"
	userpackage "nari/service/user/proto/user"
	"nari/service/user/utils"
	"time"
)

type User struct{}

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

// Call is a single request handler called via client.Call or the generated client code
func (e *User) SentSms(ctx context.Context, req *userpackage.Request, rsp *userpackage.Response) error  {
	if model.CheckCapt(req.Imgcode,req.Uuid){
		//发送短信
		credential := common.NewCredential(
			viper.GetString("SECRETID"),
			viper.GetString("SECRETKEY"),
		)
		cpf := profile.NewClientProfile()
		cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
		client, _ := sms.NewClient(credential, "ap-nanjing", cpf)

		request := sms.NewSendSmsRequest()

		request.PhoneNumberSet = common.StringPtrs([]string{ req.Phone })
		request.SmsSdkAppId = common.StringPtr("1400576568")
		request.SignName = common.StringPtr("IT人的资料屋")
		request.TemplateId = common.StringPtr("1138616")

		rand.Seed(time.Now().UnixNano()) //播种随机种子
		smscode:=fmt.Sprintf("%06d",rand.Int31n(1000000))  //生成6位随机数
		request.TemplateParamSet = common.StringPtrs([]string{ smscode })

		response, err := client.SendSms(request)
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			fmt.Printf("An API error has returned: %s", err)
			return nil
		}
		if err != nil {
			fmt.Println(err)
			rsp.Errno=utils.RECODE_SMSERR
			rsp.Errmsg=utils.RecodeText(utils.RECODE_SMSERR)
		}
		fmt.Printf("%s", response.ToJsonString())
		model.InputSms(req.Phone,smscode)
		//短信发送成功，反馈正确信息给浏览器
		rsp.Errno=utils.RECODE_OK
		rsp.Errmsg=utils.RecodeText(utils.RECODE_OK)
		//短信发送完毕
	}else {
		rsp.Errno=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(utils.RECODE_DATAERR)
	}
	return nil
}

func (e *User) RegisterModel(ctx context.Context, reg *userpackage.RegReq, rsp *userpackage.Response) error {
	if model.CheckSms(reg.Mobile,reg.Smscode){
		m5 := md5.New()
		m5.Write([]byte(reg.Password))
		m5Pwd := hex.EncodeToString(m5.Sum(nil))
		user :=model.User{
			Mobile:reg.Mobile,
			Name:reg.Mobile,
			Password_hash:m5Pwd,
		}
		err:=model.Gdb.Create(&user).Error
		if err!=nil{
			rsp.Errno=utils.RECODE_DATAERR
			rsp.Errmsg=utils.RecodeText(utils.RECODE_DATAERR)
		}else {
			rsp.Errno=utils.RECODE_OK
			rsp.Errmsg=utils.RecodeText(utils.RECODE_OK)
		}
	}else {
		rsp.Errno=utils.RECODE_SMSERR
		rsp.Errmsg=utils.RecodeText(utils.RECODE_SMSERR)
	}
	return  nil
}

func (e *User) GetallArea(ctx context.Context, req *userpackage.Request, rsp *userpackage.Resp) error {
	conn:=model.RedisPool.Get()
	var areas []model.Area
	areaData, _ := redis.Bytes(conn.Do("get", "areaData"))
	if len(areaData)==0{
		model.Gdb.Find(&areas)
		areaData,_=json.Marshal(areas)
		_, _ = conn.Do("set", "areaData", areaData)
	}
	rsp.Errno=utils.RECODE_OK
	rsp.Errmsg=utils.RecodeText(utils.RECODE_OK)
	rsp.Data=areaData
	return nil
}

func (e *User) LoginUser(ctx context.Context, req *userpackage.LoginReq, rsp *userpackage.Response) error {
	var user model.User
	m5 := md5.New()
	m5.Write([]byte(req.Password))
	m5Pwd := hex.EncodeToString(m5.Sum(nil))

	err:=model.Gdb.Where("mobile=?",req.Mobile).Where("password_hash=?",m5Pwd).
		Select("name").Find(&user).Error
	//fmt.Println(user.Name)
	if err!=nil{
		rsp.Errno=utils.RECODE_LOGINERR
		rsp.Errmsg=utils.RecodeText(utils.RECODE_LOGINERR)
	}else {
		rsp.Errno=utils.RECODE_OK
		rsp.Errmsg=utils.RecodeText(utils.RECODE_OK)
	}
	return nil
}

func (e *User) GetUserInfo(ctx context.Context, req *userpackage.UserName,rsp *userpackage.Resp) error {
	userData,err:=model.GetUserInfo(req.Username)
	if err!=nil{
		rsp.Errno=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(utils.RECODE_DATAERR)
	}else {
		rsp.Errno=utils.RECODE_OK
		rsp.Errmsg=utils.RecodeText(utils.RECODE_OK)
		rsp.Data,_=json.Marshal(userData)
	}
	return nil
}

func (e *User) UpdataUser(ctx context.Context,req *userpackage.UpdataName,rsp *userpackage.Response) error {
	err:=model.UpdataUser(req.Oldname,req.Newname)
	if err!=nil{
		rsp.Errno=utils.RECODE_DATAERR
		rsp.Errmsg=utils.RecodeText(utils.RECODE_DATAERR)
	}else {
		rsp.Errno=utils.RECODE_OK
		rsp.Errmsg=utils.RecodeText(utils.RECODE_OK)
	}
	return nil
}



