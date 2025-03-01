package tools

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/agent/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/agent/biz/dao"
	chat_models "github.com/bitdance-panic/gobuy/app/services/agent/models"
	agentutils "github.com/bitdance-panic/gobuy/app/services/agent/utils"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

var productTemplate *prompt.DefaultChatTemplate
var productSqlGenerator *openai.ChatModel

func NewSearchProductsTool() tool.BaseTool {
	InitSearchProducts()
	searchProductsTool, err := utils.InferTool("outer_search_products", "Search for products based on user requirements.", searchProductsFunc)
	if err != nil {
		log.Fatal(err)
	}
	return searchProductsTool
}

type SearchProductsParams struct {
	Prompt string `json:"sql" jsonschema:"description=User prompt"`
}

func searchProductsFunc(ctx context.Context, params *SearchProductsParams) (*ToolResponse, error) {
	log.Printf("大模型调用这个工具，prompt为: %+v", params.Prompt)
	messages, err := productTemplate.Format(context.Background(), map[string]any{
		"task": params.Prompt,
	})
	if err != nil {
		return nil, err
	}
	sqlResult, err := productSqlGenerator.Generate(ctx, messages)
	// 一般是没找到工具就进这里
	if err != nil {
		log.Fatalf("Error occurred: %v", err)
		return nil, errors.New("Tools not found")
	}
	log.Println("大模型生成的SQL:", sqlResult.Content)
	sql := agentutils.CleanBlock(sqlResult.Content)
	log.Println("clean后的SQL:", sql)
	products := make([]models.Product, 0)
	// 具体的调用逻辑
	result := tidb.DB.Raw(sql).Scan(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	var idsString strings.Builder
	// 多次追加字符串
	for i, product := range products {
		if i != 0 {
			idsString.WriteString(", ")
		}
		idsString.WriteString(strconv.Itoa(product.ID))
	}
	return &ToolResponse{
		Data:            idsString.String(),
		DataDescription: "获取的为各个商品ID,用逗号分隔",
		ShowWay:         "将各个商品ID改为超链接,格式为 http://localhost:8080/product/:id",
	}, nil
}

func InitSearchProducts() {
	columns, err := dao.GetColumns(tidb.DB)
	if err != nil {
		log.Fatal(err.Error())
	}
	columnsString := strings.Join(columns, ", ")
	log.Printf("product表字段为: %+v", columnsString)
	productSqlGenerator, productTemplate, err = chat_models.NewSearchProductSQLGenerator(context.Background(), columnsString)
	if err != nil {
		log.Panic(err.Error())
	}

}
