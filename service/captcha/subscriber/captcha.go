package subscriber

import (
	"context"
	"github.com/micro/go-micro/util/log"

	captcha "service/captcha/proto/captcha"
)

type Captcha struct{}

func (e *Captcha) Handle(ctx context.Context, msg *captcha.Message) error {
	log.Log("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *captcha.Message) error {
	log.Log("Function Received message: ", msg.Say)
	return nil
}
