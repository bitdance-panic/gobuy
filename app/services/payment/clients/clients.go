package clients

import (
	"log"

	"github.com/bitdance-panic/gobuy/app/clients"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order/orderservice"
	"github.com/bitdance-panic/gobuy/app/services/payment/conf"
	"github.com/smartwalle/alipay/v3"
)

var (
	AliPayClient *alipay.Client
	OrderClient  orderservice.Client
)

func init() {
	var err error
	AliPayClient, err = getClient()
	if err != nil {
		panic(err)
	}
	addr := conf.GetConf().Registry.RegistryAddress[0]
	OrderClient = clients.NewOrderClient(addr)
}

func getClient() (*alipay.Client, error) {
	alipayConf := conf.GetConf().Alipay
	var client *alipay.Client
	var err error
	// 支付宝提供了用于开发时测试的 sandbox 环境，对接的时候需要注意相关的 app id 和密钥是 sandbox 环境还是 production 环境的。如果是 sandbox 环境，本参数应该传 false，否则为 true。
	if client, err = alipay.New(alipayConf.APPID, alipayConf.PrivateKey, false); err != nil {
		log.Println("APPID or PrivateKey")
		return nil, err
	}
	// 加载证书
	// 加载应用公钥证书
	if err = client.LoadAppCertPublicKeyFromFile("conf/appPublicCert.crt"); err != nil {
		log.Println("appPublicCert")
		return nil, err
	}
	// 加载支付宝根证书
	if err = client.LoadAliPayRootCertFromFile("conf/alipayRootCert.crt"); err != nil {
		log.Println("alipayRootCert")
		return nil, err
	}
	// 加载支付宝公钥证书
	if err = client.LoadAlipayCertPublicKeyFromFile("conf/alipayPublicCert.crt"); err != nil {
		log.Println("alipayPublicCert")
		return nil, err
	}
	return client, nil
}
