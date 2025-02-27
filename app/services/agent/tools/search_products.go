package tools

import (
	"context"
	"log"
	"strings"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/agent/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/agent/biz/dao"
	"github.com/bitdance-panic/gobuy/app/services/agent/conf"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

var searchProductsAgent compose.Runnable[map[string]any, *schema.Message]

type searchProductsParams struct {
	SQL string `json:"sql" jsonschema:"description=SQL for searching products."`
}

type OuterSearchProductsParams struct {
	Prompt string `json:"sql" jsonschema:"description=User prompt"`
}

type SearchProductsResponse struct {
	Success  bool             `json:"success"`
	Products []models.Product `json:"products"`
}

func NewOuterSearchProductsTool() tool.BaseTool {
	searchProductsTool, err := utils.InferTool("outer_search_products", "Search for products based on user requirements.", searchOuterProductsFunc)
	if err != nil {
		log.Fatal(err)
	}
	return searchProductsTool
}

func newSearchProductsTool() tool.BaseTool {
	searchProductsTool, err := utils.InferTool("search_products", "Search for products based on user requirements.", searchProductsFunc)
	if err != nil {
		log.Fatal(err)
	}
	return searchProductsTool
}

func InitSearchProductsAgent() {
	ctx := context.Background()
	columns, err := dao.GetColumns(tidb.DB)
	if err != nil {
		log.Fatal(err)
	}
	columnsString := strings.Join(columns, ", ")
	template := prompt.FromMessages(schema.FString,
		&schema.Message{
			Role:    schema.System,
			Content: "你是一个SQL专家。现在有一个MySQL的product表,列名为：" + columnsString,
		},
		&schema.Message{
			Role:    schema.User,
			Content: "{task}。",
		},
	)
	tools := []tool.BaseTool{
		newSearchProductsTool(),
	}
	temp := float32(0.7)
	llmConf := conf.GetConf().LLM
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL:     llmConf.BaseURL,
		Model:       llmConf.ModelName,
		APIKey:      llmConf.ApiKey,
		Temperature: &temp,
	})
	if err != nil {
		log.Fatal(err)
	}
	toolInfos := make([]*schema.ToolInfo, 0, len(tools))
	for _, tool := range tools {
		info, err := tool.Info(ctx)
		if err != nil {
			log.Fatal(err)
		}
		toolInfos = append(toolInfos, info)
	}
	err = chatModel.BindTools(toolInfos)
	if err != nil {
		log.Fatal(err)
	}
	toolsNode, err := compose.NewToolNode(context.Background(), &compose.ToolsNodeConfig{
		Tools: tools,
	})
	if err != nil {
		log.Fatal(err)
	}
	chain := compose.NewChain[map[string]any, *schema.Message]()
	chain.
		AppendChatTemplate(template).
		AppendChatModel(chatModel).
		AppendToolsNode(toolsNode).
		AppendLambda(compose.InvokableLambda(func(_ context.Context, results []*schema.Message) (*schema.Message, error) {
			return results[0], nil
		}))
	searchProductsAgent, err = chain.Compile(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func searchOuterProductsFunc(ctx context.Context, params *OuterSearchProductsParams) (string, error) {
	log.Printf("大模型调用这个工具，prompt为: %+v", params.Prompt)
	variables := map[string]any{
		"task": params.Prompt,
	}
	// 调用大模型
	resp, err := searchProductsAgent.Invoke(ctx, variables)
	// 一般是没找到工具就进这里
	if err != nil {
		log.Fatalf("Error occurred: %v", err)
	}
	// 就是response的字符串
	return resp.Content, nil
}

func searchProductsFunc(_ context.Context, params *searchProductsParams) (SearchProductsResponse, error) {
	log.Printf("已生成sql,传入SearchProductFunc并准备与数据库交互: %+v", *params)
	var products []models.Product
	// 具体的调用逻辑
	result := tidb.DB.Raw(params.SQL).Scan(&products)
	if result.Error != nil {
		return SearchProductsResponse{}, result.Error
	}
	resp := SearchProductsResponse{
		Success:  true,
		Products: products,
	}
	return resp, nil
}
