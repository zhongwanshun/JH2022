package handler

import (
	"demo/db"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 登陆
func LoginHandler(c *gin.Context) {
	var loginUser db.LoginMessage
	if err := c.ShouldBind(&loginUser); err == nil { //{ID:1}
		user, err := db.Db_Login(&loginUser) // 注册用户
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		if user != nil {
			// 成功则重定向到原页面
			c.JSON(http.StatusOK, gin.H{"status": "success", "info": user})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": err.Error(), "info": "no data"})
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
	}

}

// 注册用户
func AddUsersHandler(c *gin.Context) {
	var newUser db.User
	if err := c.ShouldBind(&newUser); err == nil { //{ID:1}
		fmt.Printf("newUser:%v\n", newUser)

		err := db.Db_Register(&newUser) // 注册用户
		if err != nil {
			fmt.Printf("handler/user.go err:%v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		} else {
			// 成功则重定向到原页面
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		}

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
