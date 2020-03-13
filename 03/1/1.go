package main

import (
	"fmt"

	//导入MYSQL数据库驱动，这里使用的是GORM库封装的MYSQL驱动导入包，实际上大家看源码就知道，这里等价于导入github.com/go-sql-driver/mysql
	//这里导入包使用了 _ 前缀代表仅仅是导入这个包，但是我们在代码里面不会直接使用。
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	//配置MySQL连接参数
	username := "root"       //账号
	password := "123456"     //密码
	host := "192.168.10.188" //数据库地址，可以是Ip或者域名
	port := 3307             //数据库端口
	Dbname := "test"         //数据库名
	timeout := "10s"         //连接超时，10秒

	//拼接下dsn参数, dsn格式可以参考上面的语法，这里使用Sprintf动态拼接dsn参数，因为一般数据库连接参数，我们都是保存在配置文件里面，需要从配置文件加载参数，然后拼接dsn。
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	fmt.Println("ok.")
	//延时关闭数据库连接
	defer db.Close()
}
