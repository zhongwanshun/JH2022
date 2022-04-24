package db

// 定义数据库的访问模型
// 注册表结构体
type User struct {
	// 结构体成员tag，包含3种，一种是数据库，一种是json，一种是form表单
	ID       int    `sql:"id" json:"id" form:"id"`                   // 用户唯一标识
	Username string `sql:"username" json:"username" form:"username"` // 用户名
	Phonenum int    `sql:"phonenum" json:"phonenum" form:"phonenum"` // 电话号码
	Password string `sql:"password" json:"password" form:"password"` // 密码
	Time     string `sql:"time" json:"time" form:"time"`             //创建时间
}
type LoginMessage struct {
	Username string `form:"username" json:"username" `
	Password string `form:"password" json:"password" `
}
