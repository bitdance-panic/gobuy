package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product"
	"github.com/bitdance-panic/gobuy/app/rpc/kitex_gen/product/productservice"
	"github.com/cloudwego/kitex/client"
)

func main() {
	// 创建客户端
	client, err := productservice.NewClient("productservice", client.WithHostPorts("127.0.0.1:2680"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// 调用 CreateProduct 方法
	createReq := &product.CreateProductRequest{
		Name:        "Test Product",
		Description: "This is a test product created via client.",
		Price:       99.99,
	}
	createResp, err := client.CreateProduct(context.Background(), createReq)
	if err != nil {
		log.Fatalf("CreateProduct failed: %v", err)
	}
	fmt.Printf("CreateProduct Response: %+v\n", createResp)

	// 调用 GetProductByID 方法
	getReq := &product.GetProductByIDRequest{
		Id: createResp.Product.Id,
	}
	getResp, err := client.GetProductByID(context.Background(), getReq)
	if err != nil {
		log.Fatalf("GetProductByID failed: %v", err)
	}
	fmt.Printf("GetProductByID Response: %+v\n", getResp)

	// 调用 UpdateProduct 方法
	updateReq := &product.UpdateProductRequest{
		Id:          createResp.Product.Id,
		Name:        "Updated Test Product",
		Description: "This is an updated test product.",
		Price:       199.99,
	}
	updateResp, err := client.UpdateProduct(context.Background(), updateReq)
	if err != nil {
		log.Fatalf("UpdateProduct failed: %v", err)
	}
	fmt.Printf("UpdateProduct Response: %+v\n", updateResp)

	// 调用 DeleteProduct 方法
	deleteReq := &product.DeleteProductRequest{
		Id: createResp.Product.Id,
	}
	deleteResp, err := client.DeleteProduct(context.Background(), deleteReq)
	if err != nil {
		log.Fatalf("DeleteProduct failed: %v", err)
	}
	fmt.Printf("DeleteProduct Response: %+v\n", deleteResp)

	// 调用 SearchProducts 方法
	createReq = &product.CreateProductRequest{
		Name:        "Test666",
		Description: "This is a test product created via client.",
		Price:       99.99,
	}

	createResp, err = client.CreateProduct(context.Background(), createReq)
	if err != nil {
		log.Fatalf("Failed to create Test666 client: %v", err)
	}
	// 调用 CreateProduct 方法
	createReq = &product.CreateProductRequest{
		Name:        "666Test",
		Description: "This is a test product created 666Test via client.",
		Price:       299.99,
	}
	createResp, err = client.CreateProduct(context.Background(), createReq)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// 调用 CreateProduct 方法
	createReq = &product.CreateProductRequest{
		Name:        "777Test666",
		Description: "This is a test product created via client.",
		Price:       399.99,
	}
	createResp, err = client.CreateProduct(context.Background(), createReq)
	if err != nil {
		log.Fatalf("Failed to create 777Test666 client: %v", err)
	}

	searchReq := &product.SearchProductsRequest{
		Query: "Test",
	}
	searchResp, err := client.SearchProducts(context.Background(), searchReq)
	if err != nil {
		log.Fatalf("SearchProducts failed: %v", err)
	}
	fmt.Printf("SearchProducts Response: %+v\n", searchResp)
}
