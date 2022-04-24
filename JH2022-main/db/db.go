package db

import (
	"bytes"
	"database/sql"
	c "demo/cfg"
	"demo/log"

	_ "github.com/go-sql-driver/mysql" //init()
)

// 全局数据库对象
var DB *sql.DB

// 执行建表语句
func createTable(sql *string) error {
	_, err := DB.Exec(*sql)
	if err != nil {
		return err
	}
	return nil
}

// 初始化数据库	三张表
func InitDB(m *c.Mysql) (err error) {
	// 初始化数据库
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.User, m.Password, m.Host, m.Port, m.Database)
	var buffer bytes.Buffer
	buffer.WriteString(m.User)
	buffer.WriteString(":")
	buffer.WriteString(m.Password)
	buffer.WriteString("@tcp(")
	buffer.WriteString(m.Host)
	buffer.WriteString(":")
	buffer.WriteString(m.Port)
	buffer.WriteString(")")
	buffer.WriteString("/")
	buffer.WriteString(m.Database)
	dsn := buffer.String()

	DB, err = sql.Open("mysql", dsn) // 不会校验用户名和密码是否正确
	if err != nil {                  // dsn格式不对会报错
		log.Error.Printf("dsn:%s invalid,err:%v\n", dsn, err)
		return
	}

	err = DB.Ping() //尝试连接数据库
	if err != nil {
		log.Error.Printf("Open %s failded,err:%v\n", dsn, err)
		return
	}

	// 设置数据库连接池最大连接数
	DB.SetMaxOpenConns(10)
	// 设置数据库最大空闲数
	DB.SetMaxIdleConns(5)

	// sql语句，如果没存在库表tab_books，则新建一个
	var sqlStr = `
	CREATE TABLE IF NOT EXISTS tab_user  (
		id int(0) NOT NULL AUTO_INCREMENT COMMENT '用户唯一标识',
		username varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '用户名',
		phonenum bigint(0) NOT NULL COMMENT '电话号码',
		password varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '密码',
		time varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '创建时间',
		PRIMARY KEY (id, username, phonenum) USING BTREE,
		UNIQUE INDEX username(username) USING BTREE
	  ) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic; 
	`

	// 执行建表语句
	err = createTable(&sqlStr)
	if err != nil {
		return err
	}

	return
}
