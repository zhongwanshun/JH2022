package routers

import (
	"bytes"
	c "demo/cfg"
	h "demo/handler"
	"demo/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *c.Server) {
	switch r.SetMode {
	case "DebugMode":
		gin.SetMode(gin.DebugMode) //设置模式 ReleaseMode 生产模式,DebugMode 开发模式
	case "ReleaseMode":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	gin.ForceConsoleColor()                         //强制使用控制台颜色
	router.SetTrustedProxies([]string{r.Http_Port}) // 设置代理

	// 告诉gin框架模板文件引用的静态文件去哪里找
	router.Static("/static", "static")
	// 告诉gin框架去哪里找模板文件
	router.LoadHTMLGlob("templates/*")

	// 将请求的路径记录到日志中
	router.Use(func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/" {
			path = "/index"
		}
		log.Trace.Printf("请求的路径:%s %s\n", c.Request.Method, path)
		c.Next()
	})
	// GET请求路由组
	getRouter := router.Group("/")
	{
		/* index */
		getRouter.GET("/index", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"method": "GET",
			})
		})
		var p = map[string]bool{
			"404.html":             true,
			"about.html":           true,
			"agent.html":           true,
			"agents.html":          true,
			"compare.html":         true,
			"contact.html":         true,
			"faqs.html":            true,
			"favorites.html":       true,
			"header-default.html":  true,
			"lock-screen.html":     true,
			"login.html":           true,
			"my-properties.html":   true,
			"profile.html":         true,
			"properties.html":      true,
			"property.html":        true,
			"register.html":        true,
			"submit-property.html": true,
			"terms.html":           true,
			"toolbar-2.html":       true,
		}

		getRouter.GET("/:path", func(ctx *gin.Context) {
			path := ctx.Param("path")
			if _, OK := p[path]; OK {
				ctx.HTML(http.StatusOK, path, gin.H{})
			} else {
				ctx.HTML(http.StatusOK, "404.html", gin.H{})
			}
		})
	}
	// POST请求路由组
	postRouter := router.Group("/user")
	{
		// 登陆
		postRouter.POST("/login", h.LoginHandler)
		// 注册
		postRouter.POST("/register", h.AddUsersHandler)
	}
	var buffer bytes.Buffer
	buffer.WriteString(":")
	buffer.WriteString(r.Http_Port)
	router.Run(buffer.String())
}
