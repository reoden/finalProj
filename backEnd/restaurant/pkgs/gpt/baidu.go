package gpt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"restaurant/common"
	"restaurant/entity"
	"strings"
	"time"

	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"

	ernie "github.com/anhao/go-ernie"
	"github.com/sashabaranov/go-openai"
)

type WXYYResp struct {
	ID     string       `json:"id"`
	Result string       `json:"result"`
	Usage  openai.Usage `json:"usage"`
}

type ChatWXYY struct {
	config *ChatConfig
}

func NewWXYY(config *ChatConfig) *ChatWXYY {
	return &ChatWXYY{
		config: config,
	}
}

func (wxyy *ChatWXYY) genTurboChat(ctx context.Context, question, systemPrompt string, contextArr []string, opts ...common.Option) (*ChatResp, error) {
	API_KEY := os.Getenv("WXYYApiKey")
	SECRET_KEY := os.Getenv("WXYYSecretKey")
	client := ernie.NewDefaultClient(API_KEY, SECRET_KEY)
	messages := wxyy.genContextArr(question, contextArr)

	request := ernie.ErnieBotTurboRequest{
		Messages:    messages,
		System:      systemPrompt,
		TopP:        0.5,
		Temperature: 0.5,
	}

	log.Println("CreateErnieBotTurboChatCompletion req: ", fmt.Sprintf("%+v", request))
	res, err := client.CreateErnieBotTurboChatCompletion(context.Background(), request)
	if err != nil {
		log.Println("CreateErnieBotTurboChatCompletion error", fmt.Sprintf("%+v", err))
		return nil, err
	}
	log.Println("CreateErnieBotTurboChatCompletion resp: ", fmt.Sprintf("%+v", res))
	var resp ChatResp
	resp.Content = res.Result
	resp.Usage = Usage{
		PromptTokens:     res.Usage.PromptTokens,
		CompletionTokens: res.Usage.CompletionTokens,
		TotalTokens:      res.Usage.TotalTokens,
	}
	return &resp, nil
}

func (wxyy *ChatWXYY) genFreeChat(ctx context.Context, question, systemPrompt string, contextArr []string, opts ...common.Option) (*ChatResp, error) {
	API_KEY := os.Getenv("WXYYApiKey")
	SECRET_KEY := os.Getenv("WXYYSecretKey")
	client := ernie.NewDefaultClient(API_KEY, SECRET_KEY)
	token, err := client.GetAccessToken(ctx)
	if err != nil {
		log.Println("baidu_genFreeChat_error", fmt.Sprintf("%+v", err))
		return nil, err
	}

	for i := 1; i < 3; i++ {
		url := fmt.Sprintf("https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/ernie-speed-128k?access_token=%s", *token)
		messages := wxyy.genContextArr(question, contextArr)
		data := map[string]interface{}{
			"messages": messages,
			"system":   systemPrompt,
		}
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		log.Println("CreateErnieSpeed128kChatCompletion req: ", fmt.Sprintf("%+v", messages))
		response, err := common.MakePostRequest(url, data, headers)
		if err != nil {
			log.Println("CreateErnieSpeed128kChatCompletion error", fmt.Sprintf("%+v", err))
			if strings.Contains(err.Error(), "Open api request limit reached") {
				time.Sleep(1 * time.Second)
				continue
			} else if strings.Contains(err.Error(), "Rate limit reached for RPM") {
				return nil, fmt.Errorf("模型限频，稍后再试")
			}
			return nil, err
		}

		log.Println("CreateErnieSpeed128kChatCompletion resp: ", string(response))
		var res ernie.ErnieBotTurboResponse
		json.Unmarshal(response, &res)

		var resp ChatResp
		resp.Content = res.Result
		resp.Usage = Usage{
			PromptTokens:     res.Usage.PromptTokens,
			CompletionTokens: res.Usage.CompletionTokens,
			TotalTokens:      res.Usage.TotalTokens,
		}
		return &resp, nil
	}
	return nil, fmt.Errorf("模型限频，稍后再试")
}

func (wxyy *ChatWXYY) genBot4Chat(ctx context.Context, question, systemPrompt string, contextArr []string, opts ...common.Option) (*ChatResp, error) {
	API_KEY := os.Getenv("WXYYApiKey")
	SECRET_KEY := os.Getenv("WXYYSecretKey")
	client := ernie.NewDefaultClient(API_KEY, SECRET_KEY)
	messages := wxyy.genContextArr(question, contextArr)

	request := ernie.ErnieBot4Request{
		Messages:    messages,
		System:      systemPrompt,
		TopP:        0.5,
		Temperature: 0.5,
	}

	log.Println("CreateErnieBot4ChatCompletion req: ", fmt.Sprintf("%+v", request))
	res, err := client.CreateErnieBot4ChatCompletion(context.Background(), request)
	if err != nil {
		log.Println("CreateErnieBot4ChatCompletion error", fmt.Sprintf("%+v", err))
		return nil, err
	}
	log.Println("CreateErnieBot4ChatCompletion resp: ", fmt.Sprintf("%+v", res))
	var resp ChatResp
	resp.Content = res.Result
	resp.Usage = Usage{
		PromptTokens:     res.Usage.PromptTokens,
		CompletionTokens: res.Usage.CompletionTokens,
		TotalTokens:      res.Usage.TotalTokens,
	}
	return &resp, nil
}

func (wxyy *ChatWXYY) genTurboStreamChat(ctx context.Context, question, systemPrompt string, contextArr []string, streamChannel chan string, opts ...common.Option) error {
	API_KEY := os.Getenv("WXYYApiKey")
	SECRET_KEY := os.Getenv("WXYYSecretKey")
	client := ernie.NewDefaultClient(API_KEY, SECRET_KEY)
	messages := wxyy.genContextArr(question, contextArr)
	request := ernie.ErnieBotTurboRequest{
		Messages:    messages,
		Stream:      true,
		System:      systemPrompt,
		TopP:        0.5,
		Temperature: 0.5,
	}

	log.Println("CreateErnieBotTurboChatCompletionStream req: ", fmt.Sprintf("%+v", request))
	stream, err := client.CreateErnieBotTurboChatCompletionStream(context.Background(), request)
	if err != nil {
		log.Println("baidu_CreateErnieBotTurboChatCompletionStream_error", fmt.Sprintf("%+v", err))
		return nil
	}
	defer stream.Close()
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Println("ernie bot 4 stream EOF")
			return nil
		}
		if err != nil {
			log.Println("ernie bot 4 stream error", fmt.Sprintf("%+v", err))
			return nil
		}
		streamChannel <- response.Result

		if response.IsEnd {
			close(streamChannel)
			return nil
		}
	}
}

func (wxyy *ChatWXYY) genBot4StreamChat(ctx context.Context, question, systemPrompt string, contextArr []string, streamChannel chan string, opts ...common.Option) error {
	API_KEY := os.Getenv("WXYYApiKey")
	SECRET_KEY := os.Getenv("WXYYSecretKey")
	client := ernie.NewDefaultClient(API_KEY, SECRET_KEY)
	messages := wxyy.genContextArr(question, contextArr)
	request := ernie.ErnieBot4Request{
		Messages:    messages,
		Stream:      true,
		System:      systemPrompt,
		TopP:        0.5,
		Temperature: 0.5,
	}

	log.Println("CreateErnieBot4ChatCompletionStream req: ", fmt.Sprintf("%+v", request))
	stream, err := client.CreateErnieBot4ChatCompletionStream(context.Background(), request)
	if err != nil {
		log.Println("baidu_CreateErnieBot4ChatCompletionStream_error", fmt.Sprintf("%+v", err))
		streamChannel <- err.Error()
		close(streamChannel)
		return nil
	}
	defer stream.Close()
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			log.Println("ernie bot 4 stream EOF")
			return nil
		}
		if err != nil {
			log.Println("ernie bot 4 stream error", fmt.Sprintf("%+v", err))
			return nil
		}
		streamChannel <- response.Result

		if response.IsEnd {
			close(streamChannel)
			return nil
		}
	}
}

func (wxyy *ChatWXYY) genFreeStreamChat(ctx context.Context, question, systemPrompt string, contextArr []string, streamChannel chan string, opts ...common.Option) error {
	API_KEY := os.Getenv("WXYYApiKey")
	SECRET_KEY := os.Getenv("WXYYSecretKey")
	qianfan.GetConfig().AK = API_KEY
	qianfan.GetConfig().SK = SECRET_KEY

	// 可以通过 WithModel 指定模型

	arr := wxyy.genContextArr(question, contextArr)
	messages := make([]qianfan.ChatCompletionMessage, 0, len(arr))
	for index, m := range arr {
		if index%2 == 0 {
			messages = append(messages, qianfan.ChatCompletionUserMessage(m.Content))
		} else {
			messages = append(messages, qianfan.ChatCompletionAssistantMessage(m.Content))
		}
	}
	chat := qianfan.NewChatCompletion(
		qianfan.WithModel("ERNIE-Speed-128K"), // 支持的模型可以通过 chat.ModelList() 获取
	)

	resp := &qianfan.ModelResponseStream{}
	var err error
	for i := 1; i < 3; i++ {
		resp, err = chat.Stream( // Stream 启用流式返回，参数与 Do 相同
			context.TODO(),
			&qianfan.ChatCompletionRequest{
				Messages: messages,
				System:   systemPrompt,
			},
		)
		if err != nil {
			fmt.Println(i, resp, err)
			log.Println("baidu_ernie_speed_stream_error", fmt.Sprintf("%+v", err))
			if strings.Contains(err.Error(), "Open api request limit reached") {
				time.Sleep(1 * time.Second)
				continue
			} else if strings.Contains(err.Error(), "Rate limit reached for RPM") {
				break
			}
			continue
		}
		for {
			r, err := resp.Recv()
			if err != nil {
				return err
			}
			if resp.IsEnd { // 判断是否结束
				close(streamChannel)
				return nil
			}
			streamChannel <- r.Result
		}
	}
	if err != nil {
		streamChannel <- err.Error()
		close(streamChannel)
	}
	return nil
}

func (wxyy *ChatWXYY) genContextArr(question string, contextArr []string) []ernie.ChatCompletionMessage {
	var res []ernie.ChatCompletionMessage
	for i := len(contextArr); i > 0; {
		if contextArr[i-1] == "" || contextArr[i-2] == "" {
			i -= 2
			continue
		}
		res = append(res, ernie.ChatCompletionMessage{
			Role:    ernie.MessageRoleUser,
			Content: contextArr[i-2],
		})
		res = append(res, ernie.ChatCompletionMessage{
			Role:    ernie.MessageRoleAssistant,
			Content: contextArr[i-1],
		})
		i -= 2
	}

	res = append(res, ernie.ChatCompletionMessage{
		Role:    ernie.MessageRoleUser,
		Content: question,
	})
	return res
}

func (wxyy *ChatWXYY) Chat(ctx context.Context, question, systemPrompt string, contextArr []string, opts ...common.Option) (*ChatResp, error) {
	if wxyy.config.ChatModel == entity.WXYY {
		return wxyy.genBot4Chat(ctx, question, systemPrompt, contextArr, opts...)
	} else if wxyy.config.ChatModel == entity.WXYYFree {
		return wxyy.genFreeChat(ctx, question, systemPrompt, contextArr, opts...)
	}
	return wxyy.genTurboChat(ctx, question, systemPrompt, contextArr, opts...)
}

func (wxyy *ChatWXYY) StreamChat(ctx context.Context, question, systemPrompt string, contextArr []string, streamChannel chan string, opts ...common.Option) error {
	if wxyy.config.ChatModel == entity.WXYY {
		return wxyy.genBot4StreamChat(ctx, question, systemPrompt, contextArr, streamChannel, opts...)
	} else if wxyy.config.ChatModel == entity.WXYYFree {
		return wxyy.genFreeStreamChat(ctx, question, systemPrompt, contextArr, streamChannel, opts...)
	}
	return wxyy.genTurboStreamChat(ctx, question, systemPrompt, contextArr, streamChannel, opts...)
}
