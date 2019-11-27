package database

import (
	"SecKill/model"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

func init() {
	fmt.Println("Init mysql connections.")
	//创建一个数据库的连接
	var err error
	Db, err = gorm.Open("mysql", "root:shen6508@/ginhello?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database" + err.Error())
	}

	// 设置连接池连接数
	Db.DB().SetMaxOpenConns(10)
	Db.DB().SetMaxIdleConns(10)

	// 初始化数据库
	user := model.User{}
	coupon := &model.Coupon{}

	// 创建表
	tables := []interface{}{user, coupon}
	for _, table := range tables {
		if !Db.HasTable(table) {
			Db.AutoMigrate(table)
		}
	}

	// 创建唯一索引
	Db.Model(user).AddUniqueIndex("username_index", "username")
	Db.Model(coupon).AddUniqueIndex("coupon_index", "username", "coupon_name")


	//Db.Model(credit_card).
	//	AddForeignKey("owner_id", "users(id)", "RESTRICT", "RESTRICT").
	//	AddUniqueIndex("unique_owner", "owner_id")
}
