package initialize

import (
	"fmt"
	"log"
	"mxshop_srvs/user_srv/global"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitDB() {
	// dsn := "abc:123456@tcp(127.0.0.1:3306)/mxshop_user_srv?charset=utf8&parseTime=True&loc=Local"
	m := global.ServerConfig.MysqlInfo
	// fmt.Println(m)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", m.User, m.Password, m.Host, m.Port, m.Name)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
}
