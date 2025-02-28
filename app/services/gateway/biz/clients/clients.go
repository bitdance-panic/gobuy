package clients

import (
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart/cartservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order/orderservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product/productservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user/userservice"
	cclient "github.com/bitdance-panic/gobuy/app/services/cart/biz/clients"
	oclient "github.com/bitdance-panic/gobuy/app/services/order/biz/clients"

	// payclient "github.com/bitdance-panic/gobuy/app/services/payment/biz/clients"
	pclient "github.com/bitdance-panic/gobuy/app/services/product/biz/clients"
	uclient "github.com/bitdance-panic/gobuy/app/services/user/biz/clients"
	//"github.com/cloudwego/kitex/client"
)

var (
	UserClient    userservice.Client
	ProductClient productservice.Client
	// PaymentClient paymentservice.Client  // 暂未完成
	OrderClient orderservice.Client
	CartClient  cartservice.Client
)

func NewClients() {
	UserClient = uclient.Init()
	ProductClient = pclient.Init()
	// PaymentClient = payclient.PaymentClient  // 暂未完成
	OrderClient = oclient.Init()
	CartClient = cclient.Init()

	// var err error

	// UserClient, err = userservice.NewClient("user", client.WithHostPorts("0.0.0.0:8881"))
	// if err != nil {
	// 	hlog.Fatal(err)
	// }

	// ProductClient, err = productservice.NewClient("product", client.WithHostPorts("0.0.0.0:8882"))
	// if err != nil {
	// 	hlog.Fatal(err)
	// }

	// PaymentClient, err = paymentservice.NewClient("payment", client.WithHostPorts("0.0.0.0:8883"))
	// if err != nil {
	// 	hlog.Fatal(err)
	// }

	// OrderClient, err = orderservice.NewClient("order", client.WithHostPorts("0.0.0.0:8884"))
	// if err != nil {
	// 	hlog.Fatal(err)
	// }

	// CartClient, err = cartservice.NewClient("cart", client.WithHostPorts("0.0.0.0:8885"))
	// if err != nil {
	// 	hlog.Fatal(err)
	// }
}
