package dal

import (
	"github.com/bitdance-panic/gobuy/app/services/product/biz/dal/postgres"
)

func Init() {
	postgres.Init()
}
