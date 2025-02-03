package dal

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() {
	//dsn 配置按照实际改的，还没测试 应该是dsn: "%s:%s@tcp(%s:%s)/gobuy?tls=tidb&charset=utf8mb4"
    dsn: "%s:%s@tcp(%s:%s)/gobuy?tls=tidb&charset=utf8mb4"

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to TiDB successfully")
}
