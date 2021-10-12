module service/captcha

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/afocus/captcha v0.0.0-20191010092841-4bd1f21c8868
	github.com/dchest/captcha v0.0.0-20200903113550-03f5f0333e1f // indirect
	github.com/gin-gonic/gin v1.7.4
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins v1.5.1
	github.com/spf13/viper v1.9.0
	github.com/tencentcloud/tencentcloud-sdk-go v1.0.258
	google.golang.org/protobuf v1.27.1
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
