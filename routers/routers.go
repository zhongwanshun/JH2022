package routers

import (
	"net/http"

	"sharehouse/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	router := gin.Default()
	// 告诉gin框架模板文件引用的静态文件去哪里找
	router.Static("/static", "static")
	// 告诉gin框架去哪里找模板文件
	router.LoadHTMLGlob("templates/*")

	// 路由
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"method": "GET",
		})
	})
	router.POST("/", handler.AddUsersHandler)

	return router
}
