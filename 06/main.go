package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

//商品
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

func main() {
	// 获取当前时间
	//t := time.Now()
	// 只需要格式化年月日 必须按预定义模板 点开Format跳转的文件：time>format.go 文件73行 const变量
	//tpl1 := "2006-01-02 15:04:05"
	//fmt.Println(t.Format(tpl1))
	//tpl2 := "2006-01-02"
	//fmt.Println(t.Format(tpl2))

	//time.Now()是当前时间（Time类型）
	//fmt.Println(time.Now())
	//Time类型.Unix  是将Time类型转为时间戳
	//fmt.Println("now", time.Now().Unix())
	//time.Unix  是time包里的函数，将时间戳转为Time类型
	//fmt.Println(time.Unix(timestamp, 0))

	timestamp := time.Now().Unix()

	food := Food{
		Title:      "雪碧",
		Price:      1000,
		Stock:      22,
		Type:       1,
		CreateTime: time.Unix(timestamp, 0),
	}

	food2 := food

	db := getDb()

	split("create")

	if err := db.Create(&food).Error; err != nil {
		fmt.Println("插入失败", err)
		return
	} else {
		fmt.Println("插入一行数据")
	}

	split("update：更新单个字段值")

	//先查询一条记录, 保存在模型变量food
	//等价于: SELECT * FROM `foods`  WHERE (id = '2') LIMIT 1
	db.Where("id = ?", 2).Take(&food2)

	//修改food模型的值
	food.Price = 100

	//等价于: UPDATE `foods` SET `title` = '可乐', `type` = '0', `price` = '100', `stock` = '26', `create_time` = '2018-11-06 11:12:04'  WHERE `foods`.`id` = '2'
	db.Save(&food2)

	split("Updates:更新多个值")

	//例子1：
	//通过结构体变量设置更新字段
	updataFood := Food{
		Price: 120,
		Title: "柠檬雪碧",
	}

	//根据food模型更新数据库记录
	//等价于: UPDATE `foods` SET `price` = '120', `title` = '柠檬雪碧'  WHERE `foods`.`id` = '2'
	//Updates会忽略掉updataFood结构体变量的零值字段, 所以生成的sql语句只有price和title字段。
	db.Model(&food).Updates(&updataFood)

	//例子2:
	//根据自定义条件更新记录，而不是根据模型id
	updataFood2 := Food{
		Stock: 120,
		Title: "柠檬雪碧",
	}

	//设置Where条件，Model参数绑定一个空的模型变量
	//等价于: UPDATE `foods` SET `stock` = '120', `title` = '柠檬雪碧'  WHERE (price > '10')
	db.Model(Food{}).Where("price > ?", 10).Updates(&updataFood2)

	//例子3:
	//如果想更新所有字段值，包括零值，就是不想忽略掉空值字段怎么办？
	//使用map类型，替代上面的结构体变量

	//定义map类型，key为字符串，value为interface{}类型，方便保存任意值
	data := make(map[string]interface{})
	data["stock"] = 0 //零值字段
	data["price"] = 35

	//等价于: UPDATE `foods` SET `price` = '35', `stock` = '0'  WHERE (id = '2')
	db.Model(Food{}).Where("id = ?", 2).Updates(data)

	split("更新表达式")

	//等价于: UPDATE `foods` SET `stock` = stock + 1  WHERE `foods`.`id` = '2'
	db.Model(&food).Update("stock", gorm.Expr("stock + 1"))
}
