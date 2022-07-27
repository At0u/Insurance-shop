package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"mxshop_srvs/info_srv/model"
	"os"
	"time"
)


var tempDB *gorm.DB

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local","root","root","192.168.1.107",3306,"mxshop_info_srv")
	newLogger := logger.New(
		log.New(os.Stdout,"\r\n",log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	var err error
	tempDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}
}

func Test(){

	productId := 10

	for ;productId<=27;productId++{
		goodsInfo := model.GoodsSellCount{
			ProductId: int32(productId),
			SellNum: 0,
		}
		tempDB.Save(&goodsInfo)
	}
}


func main() {
	Init()
	Test()
	//TestCreateUser()
}
