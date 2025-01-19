package postgres

import (
	"fmt"
	"log"
	"os"

	"github.com/bitdance-panic/gobuy/app/services/user/conf"

	// "gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "gorm.io/plugin/opentelemetry/tracing"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	conf_ := conf.GetConf()
	dsn := fmt.Sprintf(conf_.Postgres.DSN, conf_.Postgres.User, conf_.Postgres.Password, conf_.Postgres.Host)
	log.Printf("dsn:%s", dsn)
	DB, err = gorm.Open(postgres.Open(dsn),
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
		// log.Println("进行user表的migrate")
		// DB.AutoMigrate(
		// 	&models.User{},
		// 	&models.Role{},
		// 	&models.Permission{},
		// 	&models.UserRole{},
		// 	&models.RolePermission{},
		// 	&models.Cart{},
		// 	&models.CartItem{},
		// 	&models.Order{},
		// 	&models.OrderItem{},
		// 	&models.Product{},
		// )
	}
}
