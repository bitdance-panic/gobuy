package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bitdance-panic/gobuy/app/services/agent/tools"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

var (
	agent     compose.Runnable[map[string]any, []*schema.Message]
	ctx       context.Context
	variables map[string]any
)

func setup() {
	ctx = context.Background()

	// 初始化 tools
	tools := []tool.BaseTool{
		tools.NewSearchProductsTool(), // 使用 InferTool 方式
	}

	// 创建并配置 ChatModel
	temp := float32(0.7)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// 读取环境变量
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	openAIBaseURL := os.Getenv("OPENAI_BASE_URL")
	openAIModelName := os.Getenv("OPENAI_MODEL_NAME")
	chatModel, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		BaseURL:     openAIBaseURL,
		Model:       openAIModelName,
		APIKey:      openAIAPIKey,
		Temperature: &temp,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 获取工具信息, 用于绑定到 ChatModel
	toolInfos := make([]*schema.ToolInfo, 0, len(tools))
	for _, tool := range tools {
		info, err := tool.Info(ctx)
		if err != nil {
			log.Fatal(err)
		}
		toolInfos = append(toolInfos, info)
	}

	// 将 tools 绑定到 ChatModel
	err = chatModel.BindTools(toolInfos)
	if err != nil {
		log.Fatal(err)
	}

	// 创建 tools 节点
	toolsNode, err := compose.NewToolNode(context.Background(), &compose.ToolsNodeConfig{
		Tools: tools,
	})
	if err != nil {
		log.Fatal(err)
	}

	template := prompt.FromMessages(schema.FString,
		&schema.Message{
			Role:    schema.System,
			Content: "你是一个{role}。",
		},
		&schema.Message{
			Role:    schema.User,
			Content: "{task}。",
		},
	)
	// 构建完整的处理链
	chain := compose.NewChain[map[string]any, []*schema.Message]()
	chain.
		AppendChatTemplate(template, compose.WithNodeName("template")).
		AppendChatModel(chatModel, compose.WithNodeName("chat_model")).
		AppendToolsNode(toolsNode, compose.WithNodeName("tools"))

	// 编译并运行 chain
	agent, err = chain.Compile(ctx)
	if err != nil {
		log.Println("这里报错")
		log.Fatal(err)
	}
}

func main() {
	setup()
	// 提前准备的设置
	variables = map[string]any{
		"role": "SQL专家",
		"task": "商品搜索。product表的字段有id、name、color、price、comment_num、is_deleted。只能选择未删除的商品即is_deleted=false的。分页查询,每页10个,每次只搜第一页。",
	}
	// 用户输入
	userInput := "查询评论数大于100、价格小于50的红色手机壳"
	// userInput := "给我当天天气"
	variables["task"] = fmt.Sprintf("%s%s", variables["task"], userInput)
	// 调用大模型
	resp, err := agent.Invoke(ctx, variables)
	// 一般是没找到工具就进这里
	if err != nil {
		log.Fatalf("Error occurred: %v", err)
	}
	// 输出结果
	for _, msg := range resp {
		var respString string
		err := json.Unmarshal([]byte(msg.Content), &respString)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}
		var respData tools.SearchProductsResponse
		err = json.Unmarshal([]byte(respString), &respData)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}
		fmt.Printf("%+v", respData)
	}
}
