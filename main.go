package main

import (
	"fmt"
	"sharehouse/cfg"
	"sharehouse/db"
	"sharehouse/routers"
)

func main() {

	// 载入配置文件	./cfg.json
	cfgPath := "./cfg.json"
	c, err := cfg.LoadConfig(cfgPath)
	if err != nil {
		fmt.Printf("载入配置文件错误:%v\n", err)
		return
	}
	fmt.Printf("%#v\n", c)
	// 初始化数据库
	err = db.InitDB(c)
	if err != nil {
		fmt.Printf("初始化数据库错误:%v\n", err)
		return
	}
	// 注册路由
	router := routers.SetupRouter()

	router.Run(fmt.Sprintf("%s:%s", c.Host, c.Port))
}
