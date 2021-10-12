package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/spf13/viper"
	modelpackage "captcha/model"
	cappackage "captcha/proto/captcha"
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

type Captcha struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Captcha) Call(ctx context.Context, req *cappackage.Request, rsp *cappackage.Response) error {
	cap := captcha.New()
	cap.SetFont("./conf/comic.ttf")
	img, str := cap.Create(4,captcha.NUM)
	modelpackage.InputCapt(req.Uuid,str)
	imgBuf,_:=json.Marshal(img)
	rsp.Img=imgBuf
	return nil
}
