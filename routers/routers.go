package routers

import (
	h "demo/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	router := gin.Default()

	// 告诉gin框架模板文件引用的静态文件去哪里找
	router.Static("/static", "static")
	// 告诉gin框架去哪里找模板文件
	router.LoadHTMLGlob("templates/*")

	// GET请求路由组
	getRouter := router.Group("/")
	{
		// 路由
		getRouter.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/index")
		})

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

	return router
}
