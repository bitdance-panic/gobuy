package clients

import (
	"log"

	"github.com/bitdance-panic/gobuy/app/common/clientsuite"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/agent/agentservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/cart/cartservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/order/orderservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product/productservice"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/user/userservice"
	"github.com/cloudwego/kitex/client"
)

func NewUserClient(registryAddress string) userservice.Client {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: "user",
			RegistryAddr:       registryAddress,
		}),
	}
	client, err := userservice.NewClient(
		"user",
		opts...,
	)
	if err != nil {
		log.Fatalf("faild to new user client: %v", err)
	}
	return client
}

func NewCartClient(registryAddress string) cartservice.Client {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: "cart",
			RegistryAddr:       registryAddress,
		}),
	}
	client, err := cartservice.NewClient(
		"cart",
		opts...,
	)
	if err != nil {
		log.Fatalf("faild to new cart client: %v", err)
	}
	return client
}

func NewAgentClient(registryAddress string) agentservice.Client {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: "agent",
			RegistryAddr:       registryAddress,
		}),
	}
	client, err := agentservice.NewClient(
		"agent",
		opts...,
	)
	if err != nil {
		log.Fatalf("faild to new agent client: %v", err)
	}
	return client
}

func NewOrderClient(registryAddress string) orderservice.Client {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: "order",
			RegistryAddr:       registryAddress,
		}),
	}
	client, err := orderservice.NewClient(
		"order",
		opts...,
	)
	if err != nil {
		log.Fatalf("faild to new order client: %v", err)
	}
	return client
}

func NewProductClient(registryAddress string) productservice.Client {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: "product",
			RegistryAddr:       registryAddress,
		}),
	}
	client, err := productservice.NewClient(
		"product",
		opts...,
	)
	if err != nil {
		log.Fatalf("faild to new product client: %v", err)
	}
	return client
}
