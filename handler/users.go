package handler

import (
	"fmt"
	"net/http"
	"sharehouse/db"

	"github.com/gin-gonic/gin"
)

// 添加用户
func AddUsersHandler(c *gin.Context) {
	var newUser db.Register
	if err := c.ShouldBind(&newUser); err == nil { //{ID:1}
		fmt.Printf("newUser:%v\n", newUser)
		err := db.Adduser(&newUser)
		if err != nil {
			fmt.Printf("handler/user.go err:%v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		// 成功则重定向到原页面
		c.JSON(http.StatusOK, gin.H{"error": 0})

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
