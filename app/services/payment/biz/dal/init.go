package dal

import (
	"github.com/bitdance-panic/gobuy/app/services/payment/biz/dal/tidb"
)

func Init() {
	tidb.Init()
}
