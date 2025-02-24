package tools

import (
	"context"
	"log"
	"strings"

	"github.com/bitdance-panic/gobuy/app/models"
	"github.com/bitdance-panic/gobuy/app/services/agent/biz/conf"
	"github.com/bitdance-panic/gobuy/app/services/agent/biz/dal/tidb"
	"github.com/bitdance-panic/gobuy/app/services/agent/biz/dao"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

var searchOrdersAgent compose.Runnable[map[string]any, *schema.Message]

type searchOrdersParams struct {
	SQL string `json:"sql" jsonschema:"description=SQL for searching orders."`
}

type OuterSearchOrdersParams struct {
	Prompt string `json:"sql" jsonschema:"description=User prompt"`
}

type SearchOrdersResponse struct {
	Success bool           `json:"success"`
	Orders  []models.Order `json:"orders"`
}

func NewOuterSearchOrdersTool() tool.BaseTool {
	searchProductsTool, err := utils.InferTool("outer_search_orders", "Search for orders based on user requirements.", searchOuterOrdersFunc)
	if err != nil {
		log.Fatal(err)
	}
	return searchProductsTool
}

func newSearchOrdersTool() tool.BaseTool {
	searchProductsTool, err := utils.InferTool("search_orders", "Search for orders based on user requirements.", searchOrdersFunc)
	if err != nil {
		log.Fatal(err)
	}
	return searchProductsTool
}

func init() {
	ctx := context.Background()
	columns, err := dao.GetColumns(tidb.DB)
	columnsString := strings.Join(columns, ", ")
	template := prompt.FromMessages(schema.FString,
		&schema.Message{
			Role:    schema.System,
			Content: "你是一个SQL专家。现在有一个MySQL的order表,列名为：" + columnsString,
		},
		&schema.Message{
			Role:    schema.User,
			Content: "{task}。",
		},
	)
	tools := []tool.BaseTool{
		newSearchOrdersTool(),
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
	searchOrdersAgent, err = chain.Compile(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func searchOuterOrdersFunc(ctx context.Context, params *OuterSearchOrdersParams) (string, error) {
	log.Printf("大模型调用这个工具，prompt为: %+v", params.Prompt)
	variables := map[string]any{
		"task": params.Prompt,
	}
	// 调用大模型
	resp, err := searchOrdersAgent.Invoke(ctx, variables)
	// 一般是没找到工具就进这里
	if err != nil {
		log.Fatalf("Error occurred: %v", err)
	}
	// 就是response的字符串
	return resp.Content, nil
}

// 你是一个专业的MySQL数据库查询生成器。你的任务是根据用户的输入生成对应的 SQL 查询语句。用户会提供一个关于 `order` 表的查询需求，你需要根据这个需求生成一个精确的 SQL 查询语句。`order` 表的字段包括：id, created_at, updated_at, is_deleted, user_id, order_number, total_amount, status。请确保生成的 SQL 查询语句是正确的，并且只返回查询语句本身，不要添加任何额外的内容或解释。

func searchOrdersFunc(_ context.Context, params *searchOrdersParams) (SearchOrdersResponse, error) {
	log.Printf("已生成sql,传入SearchProductFunc并准备与数据库交互: %+v", *params)
	var orders []models.Order
	// 具体的调用逻辑
	result := tidb.DB.Raw(params.SQL).Scan(&orders)
	if result.Error != nil {
		return SearchOrdersResponse{}, result.Error
	}
	resp := SearchOrdersResponse{
		Success: true,
		Orders:  orders,
	}
	return resp, nil
}
