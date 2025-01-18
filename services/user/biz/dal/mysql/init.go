package mysql

import (
	"fmt"
	"log"
	"os"

	"user/biz/model"
	"user/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "gorm.io/plugin/opentelemetry/tracing"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	conf_ := conf.GetConf()
	dsn := fmt.Sprintf(conf_.MySQL.DSN, conf_.MySQL.User, conf_.MySQL.Password, conf_.MySQL.Host)
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}
	// if err := DB.Use(tracing.NewPlugin(tracing.WithoutMetrics(), tracing.WithTracerProvider(mtl.TracerProvider))); err != nil {
	// 	panic(err)
	// }
	if os.Getenv("GO_ENV") != "production" {
		//nolint:errcheck
		log.Println("进行user表的migrate")
		DB.AutoMigrate(
			&model.User{},
		)
	}
}
