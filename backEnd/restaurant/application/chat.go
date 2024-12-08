package application

import (
	"context"
	"fmt"
	"restaurant/common"
	"restaurant/entity"
	"restaurant/pkgs/gpt"
)

type ChatService struct {
	chatModel map[entity.ChatModel]gpt.IChat
}

func NewChatService() *ChatService {
	return &ChatService{
		chatModel: map[entity.ChatModel]gpt.IChat{
			entity.WXYY:      gpt.NewWXYY(&gpt.ChatConfig{ChatModel: entity.WXYY}),
			entity.WXYYFree:  gpt.NewWXYY(&gpt.ChatConfig{ChatModel: entity.WXYYFree}),
			entity.WXYYTurbo: gpt.NewWXYY(&gpt.ChatConfig{ChatModel: entity.WXYYTurbo}),
		},
	}
}

func (cs *ChatService) Chat(ctx context.Context, params *entity.ChatData, contextArr []string) (*gpt.ChatResp, error) {
	if _, ok := cs.chatModel[params.ChatModel]; !ok {
		return nil, fmt.Errorf("Chat_ChatModel: %s not exits", params.ChatModel)
	}

	res, err := cs.chatModel[params.ChatModel].Chat(ctx, params.Question, params.SystemPrompt, contextArr,
		common.WithConfig("api_key", params.ApiKey))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (as *ChatService) GetParams(word string) *entity.ChatData {
	return &entity.ChatData{
		Question:     fmt.Sprintf("将中文翻译的英文: %s", word),
		SystemPrompt: fmt.Sprintf("只需要将中文部分翻译成英文, 不用给出注释"),
		ChatModel:    entity.WXYYFree,
	}
}
