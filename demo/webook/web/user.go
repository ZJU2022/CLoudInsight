package web

import (
	"CloudInsight/demo/webook/service"
	"regexp"
)

//定义用户接口
//注册，登录，编辑，查看用户信息

type UserHandler struct {
	//ijwt.Handler
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            service.UserService
	codeSvc        service.CodeService
}
