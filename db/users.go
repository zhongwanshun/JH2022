package db

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"sharehouse/getTime"
)

//定义一个用户已存在的错误对象
var erruserexists = errors.New("user already exists")

// 添加用户信息
func Adduser(newUser *Register) (err error) {

	var user Register
	sqlStr := `Select * from tab_register where username = ?;`
	rowObj := DB.QueryRow(sqlStr, newUser.Username) //必须传指针,sql,参数
	rowObj.Scan(&user.ID, &user.Username, &user.Phonenum, &user.Password, &user.Email, &user.RetPhonenum, &user.Time)
	fmt.Printf("%#v\n", user)

	if user.ID != 0 {
		return erruserexists
	}

	// 获取现在时间日期
	_time := getTime.GetTime()
	newUser.Time = _time
	// 密码再次加密 sha1
	pwd := newUser.Password
	h := sha1.New()
	h.Write([]byte(pwd))
	bs := h.Sum(nil)
	//将[]byte转成16进制
	newUser.Password = fmt.Sprintf("%x", bs)

	sqlStr = `INSERT INTO tab_register (username,phonenum,password,email,retPhonenum,time) VALUES (?,?,?,?,?,?);`

	ret, err := DB.Exec(sqlStr, newUser.Username, newUser.Phonenum, newUser.Password, newUser.Email, newUser.RetPhonenum, newUser.Time)

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
