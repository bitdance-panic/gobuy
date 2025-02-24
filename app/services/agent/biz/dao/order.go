package dao

import (
	"fmt"

	"gorm.io/gorm"
)

func GetColumns(db *gorm.DB) ([]string, error) {
	// 执行原生 SQL 查询获取列名
	var columns []string
	db.Raw("SHOW COLUMNS FROM `order`").Pluck("Field", &columns)
	// 打印列名
	for _, column := range columns {
		fmt.Println(column)
	}
	return columns, nil
}
