package entity

import (
	"net"
	"time"

	"github.com/sashabaranov/go-openai"
)

type Status int

const (
	StatusPending Status = iota
	StatusRefused
	StatusAccepted
	StatusDeleted
	StatusSaved
)

type Application struct {
	Id           int64     `json:"id"`
	UserId       int64     `json:"user_id"`
	Pictures     string    `json:"pictures"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	Describe     string    `json:"describe"`
	Rate         float64   `json:"rate"`
	Phone        string    `json:"phone"`
	PostCode     string    `json:"post_code"`
	PostName     string    `json:"post_name"`
	Introduction string    `json:"introduction"`
	WorkBeginAt  string    `json:"work_begin_at"`
	WorkEndAt    string    `json:"work_end_at"`
	HaveVege     int       `json:"have_vege"` //0 没有素食 1 有素食
	Status       Status    `json:"status"`    // 0 - 待审核 1-审核不通过 2-审核通过 3-已删除 4-草稿
	WorkDate     string    `json:"work_date"` // 将星期* 都json序列化后存起来
	HostmanName  string    `json:"hostman_name"`
	HostmanPics  string    `json:"hostman_pics"`
	HostmanThink string    `json:"hostman_think"`
	Hot          string    `json:"hot"`   // 招牌菜 若干个序列化后存储
	Yuyue        int64     `json:"yuyue"` // 是否可以预约
	Leixing      string    `json:"leixing"`
	Per          string    `json:"per"` // 人均
	En           string    `json:"en"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type HotEnBody struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type AppEn struct {
	Name         string      `json:"name"`
	Address      string      `json:"address"`
	Describe     string      `json:"describe"`
	PostCode     string      `json:"post_code"`
	PostName     string      `json:"post_name"`
	Introduction string      `json:"introduction"`
	HostmanName  string      `json:"hostman_name"`
	HostmanThink string      `json:"hostman_think"`
	Hot          []HotEnBody `json:"hot"` // 招牌菜 若干个序列化后存储
	Leixing      string      `json:"leixing"`
}

type AppResp struct {
	Id           int64     `json:"id"`
	UserId       int64     `json:"user_id"`
	Pictures     []string  `json:"pictures"`
	PicsUrl      []string  `json:"pics_url"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	Describe     string    `json:"describe"`
	Phone        string    `json:"phone"`
	Rate         float64   `json:"rate"`
	Introduction string    `json:"introduction"`
	PostCode     string    `json:"post_code"`
	PostName     string    `json:"post_name"`
	WorkBeginAt  string    `json:"work_begin_at"`
	WorkEndAt    string    `json:"work_end_at"`
	HaveVege     int       `json:"have_vege"` //0 没有素食 1 有素食
	Status       Status    `json:"status"`    // 0 - 待审核 1-审核不通过 2-审核通过 3-已删除 4-草稿
	WorkDate     []string  `json:"work_date"` // 将星期* 都json序列化后存起来
	HostmanName  string    `json:"hostman_name"`
	HostmanPics  []string  `json:"hostman_pics"`
	HostmanUrls  []string  `json:"hostman_urls"`
	HostmanThink string    `json:"hostman_think"`
	Hot          []HotBody `json:"hot"`   // 招牌菜 若干个序列化后存储
	Yuyue        int64     `json:"yuyue"` // 是否可以预约
	Leixing      string    `json:"leixing"`
	Per          string    `json:"per"` // 人均
	En           AppEn     `json:"en"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type HotBody struct {
	HotName    string   `json:"hot_name"`
	HotPics    []string `json:"hot_pics"`
	HotPicUrls []string `json:"hot_pic_urls"`
	HotDesc    string   `json:"hot_desc"`
}

type ApplicationListBody struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Describe string `json:"describe"`
	Status   Status `json:"status"` // 0 - 待审核 1-审核不通过 2-审核通过 3-已删除 4-草稿
	Option   int64  `json:"option"`
}

type ApplicationListResp struct {
	Total int64                 `json:"total"`
	Apps  []ApplicationListBody `json:"list"`
}

type ChatModel string

var (
	GPT3      ChatModel = "gpt3.5"
	GPT3_16k  ChatModel = "gpt3.5_16k"
	GPT4      ChatModel = "gpt4"
	WXYY      ChatModel = "wxyy"       // 4.0
	WXYYFree  ChatModel = "wxyy_free"  // 4.0
	WXYYTurbo ChatModel = "wxyy_turbo" // bot_turbo_0922
	Baichuan  ChatModel = "baichuan"
	Minimax5  ChatModel = "minimax5"
	Minimax6  ChatModel = "minimax6"
	YujuAI    ChatModel = "yujuai"
	BaiduApp  ChatModel = "baidu_app"  // baidu appbuilder
	Qwen      ChatModel = "qwen"       // 通义千问 longcontext
	QwenTurbo ChatModel = "qwen_turbo" // 通义千问 turbo
	//Llama3_8B  ChatModel = "llama3_8b"   // llama3
	DoubaoPro  ChatModel = "doubao_pro"  // Doubao-pro-128k
	DoubaoLite ChatModel = "doubao_lite" // Doubao-lite-128k
)

var modelMap = map[ChatModel]string{
	GPT3:     "gpt-3.5-turbo-0125",
	GPT3_16k: "gpt-3.5-turbo-0125",
	// GPT4:     "gpt-4-turbo-2024-04-09",
	GPT4:       "gpt-4o",
	Minimax5:   "abab5.5s-chat",
	Minimax6:   "abab6-chat",
	DoubaoPro:  "ep-20240601154745-sqm8m",
	DoubaoLite: "ep-20240601154659-tqqcp",
}

func GetTokenEncoding(model ChatModel) string {
	if encoding, ok := modelMap[model]; ok {
		return encoding
	}
	return openai.GPT3Dot5Turbo16K
}

type ChatData struct {
	JWTToken          string    `json:"token"`
	AppId             int64     `json:"app_id"`  // 应用id
	ApiKey            string    `json:"api_key"` // 应用id
	Question          string    `json:"question"`
	SystemPrompt      string    `json:"system_prompt"`
	ChatModel         ChatModel `json:"chat_model"`
	ChatType          string    `json:"chat_type"`
	LastChatId        int64     `json:"last_chat_id"`
	WindowId          int64     `json:"window_id"`
	ContextArr        []string
	AssistantId       string
	AssistantType     string
	AssistantFilesLen int
	To                string
	UserAccountId     int64
	CompanyAccount    *Account
	MasterAccountId   int64
	CompanyAccountId  int64
	IsFree            bool
	Conn              net.Conn
}
