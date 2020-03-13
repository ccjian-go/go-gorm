package main

//字段注释说明了gorm库把struct字段转换为表字段名长什么样子。
type Food struct {
	Id     int     //表字段名为：id
	Name   string  //表字段名为：name
	Price  float64 //表字段名为：price
	TypeId int     //表字段名为：type_id
	//字段定义后面使用两个反引号``包裹起来的字符串部分叫做标签定义，这个是golang的基础语法，不同的库会定义不同的标签，有不同的含义
	CreateTime int64 `gorm:"column:createtime"` //表字段名为：createtime
}

//设置表名，可以通过给Food struct类型定义 TableName函数，返回一个字符串作为表名
func (v Food) TableName() string {
	return "food"
}
