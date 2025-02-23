package clients

import (
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product/productservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user/userservice"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
)

var (
	UserClient    userservice.Client
	ProductClient productservice.Client
)

func init() {
	var err error
	UserClient, err = userservice.NewClient("user", client.WithHostPorts("0.0.0.0:8881"))
	if err != nil {
		hlog.Fatal(err)
	}

	ProductClient, err = productservice.NewClient("product", client.WithHostPorts("0.0.0.0:8882"))
	if err != nil {
		hlog.Fatal(err)
	}
}
