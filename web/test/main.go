package main

import (
	"github.com/afocus/captcha"
	"image/png"
	"net/http"
)

func main(){
	cap := captcha.New()
	// 设置字体
	cap.SetFont("comic.ttf")
	// 创建验证码 4个字符 captcha.NUM 字符模式数字类型
	// 返回验证码图像对象以及验证码字符串 后期可以对字符串进行对比 判断验证
	//img,str := cap.Create(4,captcha.NUM)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		img, str := cap.Create(6,captcha.NUM)
		png.Encode(w, img)
		println(str)
	})

	http.ListenAndServe(":8085", nil)
}
