package global

import (
	"gorm.io/gorm"
	"mxshop_srvs/info_srv/config"
)

var (
	DB *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig config.NacosConfig
)

//func init() {
//	dsn := "root:root@tcp(192.168.1.103:3306)/mxshop_info_srv?charset=utf8mb4&parseTime=True&loc=Local"
//
//	newLogger := logger.New(
//		log.New(os.Stdout,"\r\n",log.LstdFlags),
//		logger.Config{
//			SlowThreshold: time.Second,
//			LogLevel:      logger.Info,
//			Colorful:      true,
//		},
//	)
//
//	var err error
//	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//		NamingStrategy: schema.NamingStrategy{
//			SingularTable: true,
//		},
//		Logger: newLogger,
//	})
//
//	if err != nil {
//		panic(err)
//	}
//
//}
