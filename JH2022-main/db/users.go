package db

import (
	"crypto/sha1"
	"demo/getTime"
	"fmt"
)

// 登陆
func Db_Login(userMessage *LoginMessage) (*User, error) {
	var username = userMessage.Username
	var password = userMessage.Password
	// 查询用户是否存在
	var user User
	sqlStr := `Select * from tab_user where username = ? and password = ?;`
	rowObj := DB.QueryRow(sqlStr, username, passwd(password)) //必须传指针,sql,参数
	err := rowObj.Scan(&user.ID, &user.Username, &user.Phonenum, &user.Password, &user.Time)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// 注册用户
func Db_Register(newUser *User) (err error) {
	var user User
	sqlStr := `Select * from tab_user where username = ?;`
	rowObj := DB.QueryRow(sqlStr, newUser.Username) //必须传指针,sql,参数
	rowObj.Scan(&user.ID, &user.Username, &user.Phonenum, &user.Password, &user.Time)

	// 获取现在时间日期
	_time := getTime.GetTime()
	newUser.Time = _time
	// 密码再次加密 sha1
	pwd := newUser.Password
	newUser.Password = passwd(pwd)
	sqlStr = `INSERT INTO tab_user (username,phonenum,password,time) VALUES (?,?,?,?);`

	ret, err := DB.Exec(sqlStr, newUser.Username, newUser.Phonenum, newUser.Password, newUser.Time)

	if err != nil {
		return
	}
	n, err := ret.RowsAffected()
	if err != nil {
		return
	}
	fmt.Printf("添加了%d行数据\n", n)
	return nil
}

// 密码加密
func passwd(pwd string) string {
	h := sha1.New()
	h.Write([]byte(pwd))
	bs := h.Sum(nil)
	//将[]byte转成16进制
	return fmt.Sprintf("%x", bs)
}
