package tools

import (
	"context"
	"encoding/json"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type SearchProductsParams struct {
	SQL string `json:"sql" jsonschema:"description=SQL for searching products."`
}

// type SearchProductsResponse struct {
// 	Success  bool             `json:"success"`
// 	Message  string           `json:"message"`
// 	Products []map[string]any `json:"products"`
// }

func NewSearchProductsTool() tool.BaseTool {
	searchProductsTool, err := utils.InferTool("search_products", "Search for products based on user requirements.", searchProductsFunc)
	if err != nil {
		log.Fatal(err)
	}
	return searchProductsTool
}

func searchProductsFunc(ctx context.Context, params *SearchProductsParams) (string, error) {
	log.Printf("已生成sql,传入SearchProductFunc并准备与数据库交互: %+v", *params)
	// 具体的调用逻辑
	// gorm.find(params.SQL, &products)
	products := []map[string]any{
		{"id": 1, "name": "手机壳1", "color": "红色", "price": 50, "comment_num": 120, "is_deleted": false},
		{"id": 2, "name": "2手机壳", "color": "红色", "price": 40, "comment_num": 130, "is_deleted": false},
	}
	resp := SearchProductsResponse{
		Success:  true,
		Message:  "查询成功",
		Products: products,
	}
	respString, _ := json.Marshal(resp)
	return string(respString), nil
}
