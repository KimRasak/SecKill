package database

import (
	"SecKill/conf"
	"SecKill/model"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

func init() {
	fmt.Println("Load mysql config.")
	config, err := conf.GetAppConfig()
	if err != nil {
		panic("failed to load database config: " + err.Error())
	}
	dbType := config.App.Database.Type
	usr := config.App.Database.User
	pwd := config.App.Database.Password
	address := config.App.Database.Address
	dbName := config.App.Database.DbName
	dbLink := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		usr, pwd, address, dbName)

	//创建一个数据库的连接
	fmt.Println("Init mysql connections.")
	Db, err = gorm.Open(dbType, dbLink)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// 设置连接池连接数
	Db.DB().SetMaxOpenConns(config.App.Database.MaxOpen)
	Db.DB().SetMaxIdleConns(config.App.Database.MaxIdle)

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
