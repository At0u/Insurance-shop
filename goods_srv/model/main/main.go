package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"mxshop_srvs/goods_srv/model"
	"os"
	"time"
)


func main() {
	dsn := "root:root@tcp(192.168.1.107:3306)/mxshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout,"\r\n",log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
		)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}
	//_ = db.AutoMigrate(&model.Category{},
	//	&model.Banner{},
	//	&model.Goods{})

	cate := model.Category{
		Name:"宠物险",
	}
	db.Create(&cate)


}
