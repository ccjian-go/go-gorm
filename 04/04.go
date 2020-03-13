package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

//定义User模型，绑定users表，ORM库操作数据库，需要定义一个struct类型和MYSQL表进行绑定或者叫映射，struct字段和MYSQL表字段一一对应
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
	//配置MySQL连接参数
	username := "root"       //账号
	password := "123456"     //密码
	host := "192.168.10.188" //数据库地址，可以是Ip或者域名
	port := 3307             //数据库端口
	Dbname := "test"         //数据库名

	//通过前面的数据库参数，拼接MYSQL DSN， 其实就是数据库连接串（数据源名称）
	//MYSQL dsn格式： {username}:{password}@tcp({host}:{port})/{Dbname}?charset=utf8&parseTime=True&loc=Local
	//类似{username}使用花括号包着的名字都是需要替换的参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
	//连接MYSQL
	db, err := gorm.Open("mysql", dsn)

	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	// 打印sql
	db.LogMode(true)

	//定义一个用户，并初始化数据
	u := User{
		Username:   "tizi365",
		Password:   "123456",
		CreateTime: time.Now().Unix(),
	}

	//插入一条用户数据
	//下面代码会自动生成SQL语句：INSERT INTO `users` (`username`,`password`,`createtime`) VALUES ('tizi365','123456','1540824823')
	//db.Create(u)

	//一般项目中我们会类似下面的写法，通过Error对象检测，插入数据有没有成功，如果没有错误那就是数据写入成功了。
	if err := db.Create(u).Error; err != nil {
		fmt.Println("插入失败", err)
		return
	}

	//获取插入记录的Id
	var id []int
	db.Raw("select LAST_INSERT_ID() as id").Pluck("id", &id)

	//因为Pluck函数返回的是一列值，返回结果是slice类型，我们这里只有一个值，所以取第一个值即可。
	fmt.Println(id[0])

	fmt.Println("执行完毕")
}
