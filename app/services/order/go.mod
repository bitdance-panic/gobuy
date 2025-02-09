//module order
module github.com/bitdance-panic/gobuy/app/services/order

go 1.22.1

replace (
	github.com/bitdance-panic/gobuy/app/models => ../../models
	github.com/bitdance-panic/gobuy/app/rpc => ../../rpc
	github.com/bitdance-panic/gobuy/app/utils => ../../utils
)
