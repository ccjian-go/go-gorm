package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

func split(str string) {
	println("=============================================================================================")
	println(str)
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

//商品模型
type Food struct {
	Id    int
	Title string
	Price float32
	Stock int
	Type  int
	//mysql datetime, date类型字段，可以和golang time.Time类型绑定， 详细说明请参考：gorm连接数据库章节。
	CreateTime time.Time
}

//为Food绑定表名
func (v Food) TableName() string {
	return "foods"
}

func main() {
	db := getDb()

	split("gorm事务用法")

	// 开启事务
	tx := db.Begin()

	food := Food{Id:2}

	//在事务中执行数据库操作，使用的是tx变量，不是db。

	//库存减一
	//等价于: UPDATE `foods` SET `stock` = stock - 1  WHERE `foods`.`id` = '2' and stock > 0
	//RowsAffected用于返回sql执行后影响的行数
	rowsAffected := tx.Model(&food).Where("stock > 0").Update("stock", gorm.Expr("stock - 1")).RowsAffected
	if rowsAffected == 0 {
		//如果更新库存操作，返回影响行数为0，说明没有库存了，结束下单流程
		//这里回滚作用不大，因为前面没成功执行什么数据库更新操作，也没什么数据需要回滚。
		//这里就是举个例子，事务中可以执行多个sql语句，错误了可以回滚事务
		tx.Rollback()
		return
	}

	// 自行补充 保存订单

	//err := tx.Create(保存订单).Error
	//
	////保存订单失败，则回滚事务
	//if err != nil {
	//	tx.Rollback()
	//} else {
	//	tx.Commit()
	//}
}
