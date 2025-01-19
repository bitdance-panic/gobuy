package dal

import (
	"user/biz/dal/postgres"
)

func Init() {
	postgres.Init()
}
