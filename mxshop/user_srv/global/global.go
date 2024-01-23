package global

import (
	"mxshop_srvs/user_srv/config"

	ut "github.com/go-playground/universal-translator"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	Translator   ut.Translator
)

func init() {
	dsn := "abc:123456@tcp(127.0.0.1:3306)/mxshop_user_srv?charset=utf8&parseTime=True&loc=Local"
	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags),
	// 	logger.Config{
	// 		SlowThreshold: time.Second,
	// 		LogLevel:      logger.Info,
	// 		Colorful:      true,
	// 	},
	// )
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
}
