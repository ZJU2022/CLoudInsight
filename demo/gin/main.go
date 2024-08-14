package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求日志
		c.Next() // 继续处理请求
	}
}

// 身份验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "valid-token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort() // 终止请求
			return
		}
		c.Next() // 继续处理请求
	}
}

func main() {
	//初始化Gin实例
	r := gin.Default()

	// 注册中间件
	r.Use(LoggerMiddleware())
	r.Use(AuthMiddleware())

	//定义路由
	//静态路由
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	// 参数路由，路径参数
	r.GET("/users/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		ctx.String(http.StatusOK, "hello, "+name)
	})

	// 查询参数
	// GET /order?id=123
	r.GET("/order", func(ctx *gin.Context) {
		id := ctx.Query("id")
		ctx.String(http.StatusOK, "订单 ID 是 "+id)
	})

	r.GET("/views/*.html", func(ctx *gin.Context) {
		view := ctx.Param(".html")
		ctx.String(http.StatusOK, "view 是 "+view)
	})
	// 处理 POST 请求
	r.POST("/user", func(c *gin.Context) {
		var json struct {
			Name  string `json:"name" binding:"required"`
			Email string `json:"email" binding:"required,email"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "created",
			"name":   json.Name,
			"email":  json.Email,
		})
	})

	// 启动服务器
	r.Run(":8080")
}

/*
+------------------+      +-------------+     +--------------+
|  HTTP Request    | ---> | gin.Engine  | --> | RouterGroup  |
|  (GET /hello)    |      +-------------+     +--------------+
+------------------+             |                  |
                                 |                  |
                                 v                  v
                          +----------------------------+
                          |       Route Matching       |
                          |   (GET /hello matched)     |
                          +----------------------------+
                                 |
                                 v
						  +-------------------+
                          |  已注册Middleware |
                          +-------------------+
                                 |
                                 v

                          +-------------------+
                          |    Handler        |
                          | (Business Logic)  |
                          |  c.String(200,    |
                          |  "Hello, World!") |
                          +-------------------+
                                 |
                                 v
                          +-------------------+
                          |  HTTP Response    |
                          |  (200 OK, "Hello, |
                          |   World!")        |
                          +-------------------+
                                 |
                                 v
                          +-------------------+
                          |  Client Receives  |
                          |  Response         |
                          +-------------------+
*/
