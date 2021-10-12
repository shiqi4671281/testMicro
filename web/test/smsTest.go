package main

import (
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

func main() {

	credential := common.NewCredential(
		"AKIDTq2IMDpDMGvOlvqb8RXOzNXUCeroW3PV",
		"va0B8rLejbHhsoVxDE0WnpdEZeVAZjeK",
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, _ := sms.NewClient(credential, "ap-nanjing", cpf)

	request := sms.NewSendSmsRequest()

	request.PhoneNumberSet = common.StringPtrs([]string{ "19524270032" })
	request.SmsSdkAppId = common.StringPtr("1400576568")
	request.SignName = common.StringPtr("IT人的资料屋")
	request.TemplateId = common.StringPtr("1138616")
	request.TemplateParamSet = common.StringPtrs([]string{ "423657" })

	response, err := client.SendSms(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", response.ToJsonString())
}
