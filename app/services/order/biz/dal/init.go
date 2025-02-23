package dal

import (
	"github.com/bitdance-panic/gobuy/app/services/order/biz/dal/tidb"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	tidb.Init()
}
