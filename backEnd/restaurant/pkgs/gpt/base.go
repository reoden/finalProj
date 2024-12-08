package gpt

import (
	"context"
	"log"
	"restaurant/common"
	"restaurant/entity"

	tiktoken "github.com/pkoukk/tiktoken-go"
)

type ChatResp struct {
	Content string `json:"content"`
	Usage   Usage  `json:"usage"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatConfig struct {
	ChatModel entity.ChatModel
}

type IChat interface {
	Chat(ctx context.Context, question, systemPrompt string, contextArr []string, opts ...common.Option) (*ChatResp, error)
	StreamChat(ctx context.Context, question, systemPrompt string, contextArr []string, streamChannel chan string, opts ...common.Option) error
}

func GetChatTokens(chatModel entity.ChatModel, content string) int64 {
	model := entity.GetTokenEncoding(chatModel)
	tke, err := tiktoken.EncodingForModel(model)
	if err != nil {
		log.Printf("GetChatTokens_GetEncoding error: %v", err)
		return GetChatTokens(entity.GPT3, content)
	}

	token := tke.Encode(content, nil, nil)
	log.Printf("GetChatTokens_token: %s\n", content)
	return int64(len(token))
}
