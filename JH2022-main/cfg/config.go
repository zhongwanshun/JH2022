package cfg

import (
	"demo/log"
	"os"

	"gopkg.in/ini.v1"
)

type Server struct {
	SetMode   string `json:"setmode" ini:"setmode"`
	Protocol  string `json:"protocol" form:"protocol" ini:"protocol"`
	Http_Port string `json:"http_port" form:"http_port" ini:"http_port"`
	Pwd       string //当前目录
}
type Mysql struct {
	Host     string `sql:"host" json:"host" form:"host" ini:"host"`
	Port     string `sql:"port" json:"port" form:"port" ini:"port"`
	User     string `sql:"user" json:"user" form:"user" ini:"user"`
	Password string `sql:"password" json:"password" form:"password" ini:"password"`
	Database string `json:"database" ini:"database"`
	Clear    bool   `json:"clear" ini:"clear"`
}

type ConfigAll struct {
	Server
	Mysql
}

var config *ConfigAll

func init() {
	// 获取当前目录
	pwd, err := os.Getwd()
	if err != nil {
		log.Error.Printf("配置文件获取当前目录失败:%v", err)
	}
	// 载入配置文件
	cfg, err := ini.Load(pwd + "/config.ini")
	if err != nil {
		log.Error.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	config = &ConfigAll{
		Server: Server{
			SetMode:   cfg.Section("server").Key("setmode").In("ReleaseMode", []string{"DebugMode", "ReleaseMode"}),
			Protocol:  cfg.Section("server").Key("protocol").In("http", []string{"http", "https"}),
			Http_Port: cfg.Section("server").Key("http_port").String(),
			Pwd:       pwd,
		},
		Mysql: Mysql{
			Host:     cfg.Section("mysql").Key("host").String(),
			Port:     cfg.Section("mysql").Key("port").String(),
			User:     cfg.Section("mysql").Key("user").String(),
			Password: cfg.Section("mysql").Key("password").String(),
			Database: cfg.Section("mysql").Key("database").String(),
			Clear:    cfg.Section("mysql").Key("clear").MustBool(false), //	默认不清空数据库
		},
	}
}

func GetConfig() *ConfigAll {
	return config
}
