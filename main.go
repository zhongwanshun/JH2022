package main

import (
	c "demo/cfg"
	"demo/db"
	"demo/log"
	"demo/routers"
)

var config *c.ConfigAll

func init() {
	// 载入配置文件
	config = c.GetConfig()

	// 初始化日志库,设置日志级别
	log.NewLogger(config.Server.Pwd)

	// 初始化数据库
	mysql := config.Mysql
	err := db.InitDB(&mysql)
	if err != nil {
		log.Error.Printf("初始化数据库错误:%v\n", err)
		return
	}
}
func main() {
	// 注册路由
	server := config.Server
	routers.SetupRouter(&server)
}
