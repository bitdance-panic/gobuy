module github.com/bitdance-panic/gobuy/app/models

replace (
	github.com/bitdance-panic/gobuy/app/consts => ../../consts
	// github.com/bitdance-panic/gobuy/app/rpc => ../../rpc
	// github.com/bitdance-panic/gobuy/app/utils => ../../utils
)

go 1.22.1
