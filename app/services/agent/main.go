package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bitdance-panic/gobuy/app/services/agent/tools"
	"github.com/bitdance-panic/gobuy/app/services/user/biz/dal"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

var (
	agent compose.Runnable[map[string]any, []*schema.Message]
)

func setup() {
	ctx := context.Background()
	template := prompt.FromMessages(schema.FString,
		&schema.Message{
			Role: schema.System,
			Content: `你是一个工具挑选者，你的任务是挑选工具，并将用户的输入完全按照原样传入工具中，
			不要添加任何额外的内容或评论。不要对输入进行解释、总结或修改。`,
		},
		&schema.Message{
			Role:    schema.User,
			Content: "{task}。",
		},
	)
	tools := []tool.BaseTool{
		tools.NewSearchProductsTool(),
		tools.NewSearchOrdersTool(),
	}
	temp := float32(0.7)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
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
	chain := compose.NewChain[map[string]any, []*schema.Message]()
	chain.
		AppendChatTemplate(template).
		AppendChatModel(chatModel).
		AppendToolsNode(toolsNode)
	agent, err = chain.Compile(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	dal.Init()
	setup()
	variables := map[string]any{
		"task": "查询评论数大于100、价格小于50的红色手机壳",
	}
	// 调用大模型
	resp, err := agent.Invoke(context.Background(), variables)
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
