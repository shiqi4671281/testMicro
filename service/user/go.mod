module nari/service/user

go 1.14

require (
	github.com/gin-contrib/sessions v0.0.3 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/jinzhu/gorm v1.9.16
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins v1.5.1
	github.com/spf13/viper v1.9.0
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.258
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms v1.0.258
	google.golang.org/protobuf v1.27.1
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0