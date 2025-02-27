package models

import (
	"context"

	"github.com/bitdance-panic/gobuy/app/services/agent/conf"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func NewRespHelperModel(ctx context.Context) (*openai.ChatModel, *prompt.DefaultChatTemplate, error) {
	llmConf := conf.GetConf().LLM
	temp := float32(0.7)
	chatModel, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		BaseURL:     llmConf.BaseURL,
		Model:       llmConf.ModelName,
		APIKey:      llmConf.ApiKey,
		Temperature: &temp,
	})
	if err != nil {
		return nil, nil, err
	}
	template := prompt.FromMessages(schema.FString,
		&schema.Message{
			Role:    schema.System,
			Content: "你是一个回复整理者。你能根据用户的请求以及系统查询到的项的字符串及其描述,返回给前端HTML的项展示以及对这些项的中文总结描述。最外层不要添加```html",
		},
		&schema.Message{
			Role:    schema.User,
			Content: "用户请求为: {task}。系统查询结果为: {data},结果的描述为: {description}。展示的方式为: {showway}",
		},
	)
	return chatModel, template, nil
}
