module github.com/bitdance-panic/gobuy/app/models

replace github.com/bitdance-panic/gobuy/app/consts => ../../consts

go 1.22.1

require gorm.io/gorm v1.25.12

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.14.0 // indirect
)
