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

//定义全局的db对象，我们执行数据库操作主要通过他实现。
var _db *gorm.DB

//获取gorm db对象，其他包需要执行数据库查询的时候，只要通过tools.getDB()获取db对象即可。
//不用担心协程并发使用同样的db对象会共用同一个连接，db对象在调用他的方法的时候会从数据库连接池中获取新的连接
func GetDB() *gorm.DB {
	return _db
}

//为Food绑定表名
func (v Food) TableName() string {
	return "foods"
}

func split(str string) {
	println("=============================================================================================")
	println(str)
}

//包初始化函数，golang特性，每个包初始化的时候会自动执行init函数，这里用来初始化gorm。
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

	//设置数据库连接池参数
	db.LogMode(true)
	//	_db.DB().SetMaxOpenConns(100)   //设置数据库连接池最大连接数
	//	_db.DB().SetMaxIdleConns(20)   //连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。
	//	fmt.Println("init已执行")
	//	fmt.Println(_db)
	//}
	//
	//func main(){
	//定义接收查询结果的结构体变量
	food := Food{}

	//db := GetDB()
	//println(db)

	//查询一条记录
	//等价于：SELECT * FROM `foods`   LIMIT 1
	db.Take(&food)
	tErr := db.Take(&food).Error
	if gorm.IsRecordNotFoundError(tErr) {
		fmt.Println("查询不到数据")
	} else if tErr != nil {
		//如果err不等于record not found错误，又不等于nil，那说明sql执行失败了。
		fmt.Println("查询失败", tErr)
	}

	//查询一条记录，根据主键ID排序(正序)，返回第一条记录
	//等价于：SELECT * FROM `foods`   ORDER BY `foods`.`id` ASC LIMIT 1
	db.First(&food)

	//查询一条记录, 根据主键ID排序(倒序)，返回第一条记录
	//等价于：SELECT * FROM `foods`   ORDER BY `foods`.`id` DESC LIMIT 1
	//语义上相当于返回最后一条记录
	db.Last(&food)

	//查询多条记录，Find函数返回的是一个数组
	//因为Find返回的是数组，所以定义一个商品数组用来接收结果
	var foods []Food
	//等价于：SELECT * FROM `foods`
	db.Find(&foods)

	//查询一列值
	//商品标题数组
	var titles []string
	//返回所有商品标题
	//等价于：SELECT title FROM `foods`
	//Pluck提取了title字段，保存到titles变量
	//这里Model函数是为了绑定一个模型实例，可以从里面提取表名。
	db.Model(Food{}).Pluck("title", &titles)

	split("where")

	//例子1:
	//等价于: SELECT * FROM `foods`  WHERE (id = '10') LIMIT 1
	//这里问号(?), 在执行的时候会被10替代
	db.Where("id = ?", 10).Take(&food)

	//例子2:
	// in 语句
	//等价于: SELECT * FROM `foods`  WHERE (id in ('1','2','5','6')) LIMIT 1
	//args参数传递的是数组
	db.Where("id in (?)", []int{1, 2, 5, 6}).Take(&food)

	//例子3:
	//等价于: SELECT * FROM `foods`  WHERE (create_time >= '2018-11-06 00:00:00' and create_time <= '2018-11-06 23:59:59')
	//这里使用了两个问号(?)占位符，后面传递了两个参数替换两个问号。
	db.Where("create_time >= ? and create_time <= ?", "2018-11-06 00:00:00", "2018-11-06 23:59:59").Find(&foods)

	//例子4:
	//like语句
	//等价于: SELECT * FROM `foods`  WHERE (title like '%可乐%')
	db.Where("title like ?", "%可乐%").Find(&foods)

	split("select")

	//例子1:
	//等价于: SELECT id,title FROM `foods`  WHERE `foods`.`id` = '1' AND ((id = '1')) LIMIT 1
	db.Select("id,title").Where("id = ?", 1).Take(&food)

	//这种写法是直接往Select函数传递数组，数组元素代表需要选择的字段名
	db.Select([]string{"id", "title"}).Where("id = ?", 1).Take(&food)

	//例子2:
	//可以直接书写聚合语句
	//等价于: SELECT count(*) as total FROM `foods`
	total := []int{}

	//Model函数，用于指定绑定的模型，这里生成了一个Food{}变量。目的是从模型变量里面提取表名，Pluck函数我们没有直接传递绑定表名的结构体变量，gorm库不知道表名是什么，所以这里需要指定表名
	//Pluck函数，主要用于查询一列值
	db.Model(Food{}).Select("count(*) as total").Pluck("total", &total)

	fmt.Println(total[0])

	split("order")

	//例子:
	//等价于: SELECT * FROM `foods`  WHERE (create_time >= '2018-11-06 00:00:00') ORDER BY create_time desc
	db.Where("create_time >= ?", "2018-11-06 00:00:00").Order("create_time desc").Find(&foods)

	split("limit & offset")

	//等价于: SELECT * FROM `foods` ORDER BY create_time desc LIMIT 10 OFFSET 0
	db.Order("create_time desc").Limit(10).Offset(0).Find(&foods)

	split("count")

	//例子:
	count := 0
	//等价于: SELECT count(*) FROM `foods`
	//这里也需要通过model设置模型，让gorm可以提取模型对应的表名
	db.Model(Food{}).Count(&count)
	fmt.Println(count)

	split("group by")

	//例子:
	//统计每个商品分类下面有多少个商品
	//定一个Result结构体类型，用来保存查询结果
	type Result struct {
		Type  int
		Total int
	}

	var results []Result
	//等价于: SELECT type, count(*) as  total FROM `foods` GROUP BY type HAVING (total > 0)
	db.Model(Food{}).Select("type, count(*) as  total").Group("type").Having("total > 0").Scan(&results)

	//scan类似Find都是用于执行查询语句，然后把查询结果赋值给结构体变量，区别在于scan不会从传递进来的结构体变量提取表名.
	//这里因为我们重新定义了一个结构体用于保存结果，但是这个结构体并没有绑定foods表，所以这里只能使用scan查询函数。

	split("raw")

	//例子:
	sql := "SELECT type, count(*) as  total FROM `foods` where create_time > ? GROUP BY type HAVING (total > 0)"
	//因为sql语句使用了一个问号(?)作为绑定参数, 所以需要传递一个绑定参数(Raw第二个参数).
	//Raw函数支持绑定多个参数
	db.Raw(sql, "2018-11-06 00:00:00").Scan(&results)
	fmt.Println(results)
}
