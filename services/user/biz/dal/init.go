package dal

import (
	"user/biz/dal/mysql"
)

func Init() {
	mysql.Init()
}
