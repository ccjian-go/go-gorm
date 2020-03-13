package main

//导入tools包
import (
	"my/tools"
	"fmt"
)

//在这里User类型可以代表mysql users表
type User struct {
	//通过在字段后面的标签说明，定义golang字段和表字段的关系
	//例如 `gorm:"column:username"` 标签说明含义是: Mysql表的列名（字段名)为username
	//这里golang定义的Username变量和MYSQL表字段username一样，他们的名字可以不一样。
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	//创建时间，时间戳
	CreateTime int64 `gorm:"column:createtime"`
}

func main() {
	//获取DB
	_db := tools.GetDB()

	//执行数据库查询操作
	u := User{}
	//自动生成sql： SELECT * FROM `users`  WHERE (username = 'tizi365') LIMIT 1
	_db.Where("username = ?", "tizi365").First(&u)
}
