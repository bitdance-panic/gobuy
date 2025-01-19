package dal

import (
	"github.com/bitdance-panic/gobuy/app/services/user/biz/dal/postgres"
)

func Init() {
	postgres.Init()
}
