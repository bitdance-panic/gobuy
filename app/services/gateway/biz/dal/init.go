package dal

import (
	"github.com/bitdance-panic/gobuy/app/services/gateway/biz/dal/tidb"
)

func Init() {
	tidb.Init()
}
