package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/model"
	"os"
	"time"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local","root","root","192.168.1.107",3306,"mxshop_order_srv")
	newLogger := logger.New(
		log.New(os.Stdout,"\r\n",log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	tTime := time.Now().AddDate(1,0,0)
	oo := model.Order{
		UserId :24,

		ProductID:2,
		ProductName:"儿童医疗保险E款(互联网版)",
		CategoryID :2,
		CategoryName :"健康险",
		Term:1,

		OrderSn:"ttt",
		Status :1,
		Price :44.2,
		EndTime:&tTime,

		ApplicantName :"lzc",
		ApplicantMobile:"15159917008",
		ApplicantIdNum :"1515",

		InsurerName :"lzc",
		InsurerMobile:"15159917008",
		InsurerIdNum :"1515",
	}
	result:= global.DB.Save(&oo)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}

}
