package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

/**
查询出最后一条并删除
*/
func main() {
	db := getDb()

	split("删除模型数据")
	//======================================================================用法：db.Delete(模型变量)
	//例子：
	food := Food{}

	db.Last(&food)

	food2 := Food{}

	//先查询一条记录, 保存在模型变量food
	//等价于: SELECT * FROM `foods`  WHERE (id = '2') LIMIT 1
	db.Where("id = ?", food.Id).Take(&food2)

	//删除food对应的记录，通过主键Id标识记录
	//等价于： DELETE from `foods` where id=2;
	db.Delete(&food2)

	split("根据Where条件删除数据")

	//================================================ 用法：db.Where(条件表达式).Delete(空模型变量指针)
	//等价于：DELETE from `foods` where (`type` = 5);
	db.Where("id = ?", food.Id-1).Delete(&Food{})
}

func split(str string) {
	println("=============================================================================================")
	println(str)
}

type Food struct {
	Id    int
	Title string
	Price float32
	Stock int
	Type  int
	//mysql datetime, date类型字段，可以和golang time.Time类型绑定， 详细说明请参考：gorm连接数据库章节。
	CreateTime time.Time
}

func getDb() *gorm.DB {
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
	_db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	_db.LogMode(true)

	return _db
}
