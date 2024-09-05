package web

import (
	"CloudInsight/demo/webook/domain"
	"CloudInsight/demo/webook/service"
	"context"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

//定义用户接口
//注册，登录，编辑，查看用户信息

// UserHandler定义所有和用户有关的路由
type UserHandler struct {
	//ijwt.Handler
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            service.UserService
	codeSvc        service.CodeService
}

// NewUserHandler 构造函数，用于创建 UserHandler 实例并初始化其字段
func NewUserHandler(userSvc service.UserService, codeSvc service.CodeService) *UserHandler {
	// 初始化正则表达式用于校验邮箱和密码
	emailRexExp := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	passwordRexExp := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	return &UserHandler{
		emailRexExp:    emailRexExp,
		passwordRexExp: passwordRexExp,
		svc:            userSvc,
		codeSvc:        codeSvc,
	}
}

// 方法用于注册与用户相关的所有路由，这些路由定义了不同的 HTTP 请求与相应处理方法的映射。
func (c *UserHandler) RegisterRoutes(server *gin.Engine) {
	//分散注册
	//分组路由
	ug := server.Group("/users")
	ug.POST("/signup", c.SignUp)
	ug.POST("/login")
	ug.POST("/edit")
	ug.GET("/profile")
}

// SignUp 用户注册接口
// Web 服务接口：是指通过 HTTP 请求和响应进行通信的服务接口，例如 SignUp。
// Go 语言接口：是一种类型定义，规定了一组方法签名，由具体类型来实现。
func (h *UserHandler) SignUp(c *gin.Context) {
	//定义内部结构体来接收数据
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Name     string `json:"name" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
	}

	// 绑定请求参数并校验
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//业务逻辑

	// 检查密码复杂性
	if !h.passwordRexExp.MatchString(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码不符合要求"})
		return
	}
	//检查邮箱格式

	//密码和确认密码需要相等

	//密码需要符合一定规则

	//接口兼容

	// 创建用户结构体，将http映射成为程序中对象
	user := domain.User{
		Email:    req.Email,
		Password: req.Password, // 在实际应用中应加密存储密码
		Name:     req.Name,
		Phone:    req.Phone,
	}

	// 调用 UserService 的 Signup 方法，将新创建的 user 结构体传递给服务层，执行用户注册的业务逻辑。
	if err := h.svc.Signup(context.Background(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败，请稍后重试"})
		return
	}

	// 注册成功，返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}
