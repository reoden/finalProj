package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"restaurant/common"
	"restaurant/domain"
	"restaurant/entity"
	"restaurant/pkgs"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	officedocCommon "github.com/unidoc/unioffice/common"
	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/measurement"
)

func init() {
	err := license.SetMeteredKey(common.UNIDOC_LICENSE_API_KEY)
	if err != nil {
		log.Printf("ERROR to get license: %v", err)
		panic(err)
	}
}

type AppService struct {
	db      *domain.AppRepo
	account AccountService
	chat    ChatService
}

func NewAppService() *AppService {
	return &AppService{
		db:      domain.NewAppRepo(),
		account: *NewAccountService(),
		chat:    *NewChatService(),
	}
}

func (as *AppService) TranslateWord(ctx context.Context, param entity.AppResp) (*entity.AppEn, error) {
	nameResult, err := as.chat.Chat(ctx, as.chat.GetParams(param.Name), nil)
	if err != nil {
		log.Printf("Translate Name err: %v\n", err)
		return nil, err
	}

	addressResult, err := as.chat.Chat(ctx, as.chat.GetParams(param.Address), nil)
	if err != nil {
		log.Printf("Translate Address Error: %v\n", err)
		return nil, err
	}

	describeResult, err := as.chat.Chat(ctx, as.chat.GetParams(param.Describe), nil)
	if err != nil {
		log.Printf("Translate Describe Error: %v\n", err)
		return nil, err
	}

	postNameResult, err := as.chat.Chat(ctx, as.chat.GetParams(param.PostName), nil)
	if err != nil {
		log.Printf("Translate PostName Error: %v\n", err)
		return nil, err
	}

	introductionResult, err := as.chat.Chat(ctx, as.chat.GetParams(param.Introduction), nil)
	if err != nil {
		log.Printf("Translate Introduction Error: %v\n", err)
	}

	hostmanNameResult, err := as.chat.Chat(ctx, as.chat.GetParams(param.HostmanName), nil)
	if err != nil {
		log.Printf("Translate HostmanName Error: %v\n", err)
	}

	hostmanThinkResult, err := as.chat.Chat(ctx, as.chat.GetParams(param.HostmanThink), nil)
	if err != nil {
		log.Printf("Translate HostmanThink Error: %v\n", err)
	}

	leixingResult, err := as.chat.Chat(ctx, as.chat.GetParams(param.Leixing), nil)

	res := entity.AppEn{
		Name:         nameResult.Content,
		Address:      addressResult.Content,
		Describe:     describeResult.Content,
		PostCode:     param.PostCode,
		PostName:     postNameResult.Content,
		Introduction: introductionResult.Content,
		HostmanName:  hostmanNameResult.Content,
		HostmanThink: hostmanThinkResult.Content,
		Leixing:      leixingResult.Content,
	}

	hots := make([]entity.HotEnBody, 0)

	for i, body := range param.Hot {
		hotEnBodyNameResult, err := as.chat.Chat(ctx, as.chat.GetParams(body.HotName), nil)
		if err != nil {
			log.Printf("Translate index = %d Hotname Error: %v\n", i, err)
			return nil, err
		}

		hotEnBodyDescResult, err := as.chat.Chat(ctx, as.chat.GetParams(body.HotDesc), nil)
		if err != nil {
			log.Printf("Translate index = %d HotDesc Error: %v\n", i, err)
			return nil, err
		}

		hots = append(hots, entity.HotEnBody{
			Name: hotEnBodyNameResult.Content,
			Desc: hotEnBodyDescResult.Content,
		})
	}

	res.Hot = hots

	return &res, nil
}

type enApp struct {
	Id           int64            `json:"id"`
	UserId       int64            `json:"user_id"`
	Pictures     []string         `json:"pictures"`
	PicsUrl      []string         `json:"pics_url"`
	Name         string           `json:"name"`
	Address      string           `json:"address"`
	Describe     string           `json:"describe"`
	Phone        string           `json:"phone"`
	Rate         float64          `json:"rate"`
	Introduction string           `json:"introduction"`
	PostCode     string           `json:"post_code"`
	PostName     string           `json:"post_name"`
	WorkBeginAt  string           `json:"work_begin_at"`
	WorkEndAt    string           `json:"work_end_at"`
	HaveVege     int              `json:"have_vege"` // 0 没有素食 1 有素食
	Status       entity.Status    `json:"status"`    // 0 - 待审核 1-审核不通过 2-审核通过 3-已删除 4-草稿
	WorkDate     []string         `json:"work_date"` // 将星期* 都json序列化后存起来
	HostmanName  string           `json:"hostman_name"`
	HostmanPics  []string         `json:"hostman_pics"`
	HostmanUrls  []string         `json:"hostman_urls"`
	HostmanThink string           `json:"hostman_think"`
	Hot          []entity.HotBody `json:"hot"`   // 招牌菜 若干个序列化后存储
	Yuyue        int64            `json:"yuyue"` // 是否可以预约
	Leixing      string           `json:"leixing"`
	Per          string           `json:"per"` // 人均
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

func (as *AppService) zh2en(param entity.AppResp) entity.AppResp {
	app := new(enApp)
	app.Id = param.Id
	app.UserId = param.UserId
	app.PicsUrl = param.PicsUrl
	app.Pictures = param.Pictures
	app.Name = param.En.Name
	app.Address = param.En.Address
	app.Describe = param.En.Describe
	app.Phone = param.Phone
	app.Rate = param.Rate
	app.Introduction = param.En.Introduction
	app.PostCode = param.PostCode
	app.PostName = param.En.PostName
	app.WorkEndAt = param.WorkEndAt
	app.WorkBeginAt = param.WorkBeginAt
	app.HaveVege = param.HaveVege
	app.Status = param.Status
	app.WorkDate = param.WorkDate
	app.HostmanName = param.En.HostmanName
	app.HostmanPics = param.HostmanPics
	app.HostmanUrls = param.HostmanUrls
	app.HostmanThink = param.En.HostmanThink
	hots := make([]entity.HotBody, 0)

	for i := range param.Hot {
		hots = append(hots, entity.HotBody{
			HotPics:    param.Hot[i].HotPics,
			HotPicUrls: param.Hot[i].HotPicUrls,
			HotName:    param.En.Hot[i].Name,
			HotDesc:    param.En.Hot[i].Desc,
		})
	}

	app.Hot = hots
	app.Yuyue = param.Yuyue
	app.Leixing = param.Leixing
	app.Per = param.Per
	app.CreatedAt = param.CreatedAt
	app.UpdatedAt = param.UpdatedAt

	return entity.AppResp{
		Id:           app.Id,
		UserId:       app.UserId,
		Pictures:     app.Pictures,
		PicsUrl:      app.PicsUrl,
		Name:         app.Name,
		Address:      app.Address,
		Describe:     app.Describe,
		Phone:        app.Phone,
		Rate:         app.Rate,
		Introduction: app.Introduction,
		PostCode:     app.PostCode,
		PostName:     app.PostName,
		WorkBeginAt:  app.WorkBeginAt,
		WorkEndAt:    app.WorkEndAt,
		HaveVege:     app.HaveVege,
		Status:       app.Status,
		WorkDate:     app.WorkDate,
		HostmanName:  app.HostmanName,
		HostmanPics:  app.HostmanPics,
		HostmanUrls:  app.HostmanUrls,
		HostmanThink: app.HostmanThink,
		Hot:          app.Hot,
		Yuyue:        app.Yuyue,
		Leixing:      app.Leixing,
		Per:          app.Per,
		CreatedAt:    app.CreatedAt,
		UpdatedAt:    app.UpdatedAt,
	}
}

func (as *AppService) multiZh2En(apps []entity.AppResp) (params []entity.AppResp) {
	for _, app := range apps {
		params = append(params, as.zh2en(app))
	}

	return
}

func (as *AppService) GetApp(ctx context.Context, id int64) (*entity.AppResp, error) {
	app, err := as.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if app == nil {
		return nil, nil
	}

	param, err := as.application2Resp(app)
	if err != nil {
		return nil, err
	}

	return param, nil
}

func (as *AppService) SearchById(ctx context.Context, offset, limit int, id int64) (map[string]interface{}, error) {
	apps, total, err := as.db.SearchById(ctx, offset, limit, id)
	if err != nil {
		return nil, err
	}

	params, err := as.apps2Resp(apps)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total": total,
		"apps":  params,
	}, nil
}

func (as *AppService) GetApps(ctx context.Context, ids []int64) ([]*entity.Application, error) {
	app, err := as.db.GetApps(ctx, ids)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (as *AppService) GetAppsOrder(ctx context.Context, ids []int64) ([]*entity.AppResp, error) {
	apps, err := as.db.GetAppsOrder(ctx, ids)
	if err != nil {
		return nil, err
	}

	log.Println(apps)

	params := make([]*entity.AppResp, 0)
	for _, app := range apps {
		param := &entity.AppResp{
			Id:           app.Id,
			UserId:       app.UserId,
			Name:         app.Name,
			Address:      app.Address,
			Describe:     app.Describe,
			Rate:         app.Rate,
			Introduction: app.Introduction,
			Phone:        app.Phone,
			PostCode:     app.PostCode,
			PostName:     app.PostName,
			WorkBeginAt:  app.WorkBeginAt,
			WorkEndAt:    app.WorkEndAt,
			HaveVege:     app.HaveVege,
			Status:       app.Status,
			CreatedAt:    app.CreatedAt,
			UpdatedAt:    app.UpdatedAt,
			HostmanName:  app.HostmanName,
			HostmanThink: app.HostmanThink,
			Yuyue:        app.Yuyue,
			Leixing:      app.Leixing,
			Per:          app.Per,
		}

		var dates []string
		err = json.Unmarshal([]byte(app.WorkDate), &dates)
		if err != nil {
			return nil, err
		}
		param.WorkDate = dates

		bodyPics := make([]string, 0)
		var appPics []string
		err = json.Unmarshal([]byte(app.Pictures), &appPics)
		if err != nil {
			return nil, err
		}

		for _, pic := range appPics {
			signedUrl, _ := pkgs.SignedUrl(pic, false)
			bodyPics = append(bodyPics, signedUrl)
		}

		log.Println(bodyPics)

		param.PicsUrl = bodyPics
		param.Pictures = appPics

		var appEn entity.AppEn
		err = json.Unmarshal([]byte(app.En), &appEn)
		if err != nil {
			return nil, err
		}

		param.En = appEn

		params = append(params, param)
	}

	return params, nil
}

func (a *AppService) application2Resp(app *entity.Application) (*entity.AppResp, error) {
	param := &entity.AppResp{
		Id:           app.Id,
		UserId:       app.UserId,
		Name:         app.Name,
		Address:      app.Address,
		Describe:     app.Describe,
		Phone:        app.Phone,
		Rate:         app.Rate,
		Introduction: app.Introduction,
		PostCode:     app.PostCode,
		PostName:     app.PostName,
		WorkBeginAt:  app.WorkBeginAt,
		WorkEndAt:    app.WorkEndAt,
		HaveVege:     app.HaveVege,
		Status:       app.Status,
		CreatedAt:    app.CreatedAt,
		UpdatedAt:    app.UpdatedAt,
		HostmanName:  app.HostmanName,
		HostmanThink: app.HostmanThink,
		Yuyue:        app.Yuyue,
		Leixing:      app.Leixing,
		Per:          app.Per,
	}

	var dates []string
	err := json.Unmarshal([]byte(app.WorkDate), &dates)
	if err != nil {
		return nil, err
	}
	param.WorkDate = dates

	appPairStrs, err := a.picStr2Urls(app.Pictures)
	if err != nil {
		return nil, err
	}

	param.Pictures = appPairStrs.Pics
	param.PicsUrl = appPairStrs.Urls

	var hot []hotBody
	err = json.Unmarshal([]byte(app.Hot), &hot)
	if err != nil {
		return nil, err
	}

	hotParam := make([]entity.HotBody, 0)
	for _, hotbody := range hot {
		hotPictures, err := json.Marshal(hotbody.HotPics)
		if err != nil {
			return nil, err
		}
		hotPairStrs, err := a.picStr2Urls(string(hotPictures))
		if err != nil {
			return nil, err
		}

		hotParam = append(hotParam, entity.HotBody{
			HotName:    hotbody.HotName,
			HotDesc:    hotbody.HotDesc,
			HotPics:    hotPairStrs.Pics,
			HotPicUrls: hotPairStrs.Urls,
		})
	}

	param.Hot = hotParam

	hostmanPictures, err := a.picStr2Urls(app.HostmanPics)
	if err != nil {
		return nil, err
	}

	param.HostmanPics = hostmanPictures.Pics
	param.HostmanUrls = hostmanPictures.Urls

	var enBody entity.AppEn
	err = json.Unmarshal([]byte(app.En), &enBody)
	if err != nil {
		return nil, err
	}

	param.En = enBody

	return param, nil
}

type hotBody struct {
	HotName string   `json:"hot_name"`
	HotPics []string `json:"hot_pics"`
	HotDesc string   `json:"hot_desc"`
}

type pStr struct {
	Pics []string `json:"pics"`
	Urls []string `json:"urls"`
}

func (a *AppService) picStr2Urls(pictures string) (*pStr, error) {
	bodyPics := make([]string, 0)
	var appPics []string
	err := json.Unmarshal([]byte(pictures), &appPics)
	if err != nil {
		return nil, err
	}

	for _, pic := range appPics {
		signedUrl, _ := pkgs.SignedUrl(pic, false)
		bodyPics = append(bodyPics, signedUrl)
	}

	log.Println(bodyPics)

	return &pStr{
		Pics: appPics,
		Urls: bodyPics,
	}, nil
}

func (a *AppService) appResp2Application(app *entity.AppResp) (*entity.Application, error) {
	param := &entity.Application{
		Id:           app.Id,
		UserId:       app.UserId,
		Name:         app.Name,
		Address:      app.Address,
		Describe:     app.Describe,
		Rate:         app.Rate,
		Phone:        app.Phone,
		Introduction: app.Introduction,
		PostCode:     app.PostCode,
		PostName:     app.PostName,
		WorkBeginAt:  app.WorkBeginAt,
		WorkEndAt:    app.WorkEndAt,
		HaveVege:     app.HaveVege,
		Status:       app.Status,
		CreatedAt:    app.CreatedAt,
		UpdatedAt:    app.UpdatedAt,
		HostmanName:  app.HostmanName,
		HostmanThink: app.HostmanThink,
		Yuyue:        app.Yuyue,
		Leixing:      app.Leixing,
		Per:          app.Per,
	}

	dates, err := json.Marshal(app.WorkDate)
	if err != nil {
		return nil, err
	}

	param.WorkDate = string(dates)

	pics, err := json.Marshal(app.Pictures)
	if err != nil {
		return nil, err
	}

	param.Pictures = string(pics)

	hostmanPictures, err := json.Marshal(app.HostmanPics)
	if err != nil {
		return nil, err
	}

	param.HostmanPics = string(hostmanPictures)

	hotParams := make([]hotBody, 0)
	n := len(app.Hot)

	for i := 0; i < n; i++ {
		hotParams = append(hotParams, hotBody{
			HotName: app.Hot[i].HotName,
			HotPics: app.Hot[i].HotPics,
			HotDesc: app.Hot[i].HotDesc,
		})
	}
	hot, err := json.Marshal(hotParams)
	if err != nil {
		return nil, err
	}

	param.Hot = string(hot)

	en, err := json.Marshal(app.En)
	if err != nil {
		return nil, err
	}

	param.En = string(en)

	return param, nil
}
func sceneExtract(textContent, correctStr string) bool {
	data := map[string]string{
		"text": textContent,
	}
	req, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return false
	}

	url := "http://127.0.0.1:8000/extract"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(req))
	if err != nil {
		log.Printf("发送POST请求出错,错误信息: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return false
	}

	type respData struct {
		Scenes []string `json:"scenes"`
		Text   string   `json:"text"`
	}

	var ans respData
	if err = json.Unmarshal(body, &ans); err != nil {
		log.Println(err)
		return false
	}

	log.Println(ans.Scenes)
	ret := ans.Scenes

	for _, nowStr := range ret {
		if strings.HasPrefix(correctStr, nowStr) {
			return true
		}
	}

	return false
}

func scanMGword(textContent string) bool {
	data := map[string]string{
		"text": textContent,
	}
	req, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return true
	}

	url := "http://127.0.0.1:8000/scanMGword"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(req))
	if err != nil {
		log.Printf("发送POST请求失败,错误信息: %v\n", err)
		return true
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return true
	}

	type respData struct {
		Result int `json:"result"`
	}

	var ans respData
	if err = json.Unmarshal(body, &ans); err != nil {
		log.Println(err)
		return true
	}

	if ans.Result == 0 {
		return true
	}

	return false
}

func (a *AppService) reqPicValid(ctx context.Context, app *entity.Application) (bool, error) {
	appPairStrs, err := a.picStr2Urls(app.Pictures)
	if err != nil {
		return false, err
	}

	urls := appPairStrs.Urls
	log.Println(urls)
	return true, nil

	type reqData struct {
		Url string `json:"url"`
	}
	for _, url := range urls {
		data := map[string]string{
			"url": url,
		}
		req, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
			return false, err
		}

		reqUrl := ""
		resp, err := http.Post(reqUrl, "application/json", bytes.NewBuffer(req))
		if err != nil {
			log.Printf("发送POST请求失败,错误信息: %v\n", err)
			return false, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return false, err
		}

		type respData struct {
			FusionScore float64 `json:"fusion_score"`
		}

		var ans respData
		if err = json.Unmarshal(body, &ans); err != nil {
			log.Println(err)
			return false, err
		}

		if ans.FusionScore > 0 {
			return false, nil
		}
	}

	return true, nil
}

// 注意: openai 官方的 golang sdk 目前尚处于 alpha 阶段。未来可能会出现一些小的破坏性改动。

func (a *AppService) reqOpenai(ctx context.Context, param *entity.Application, name string) (bool, error) {
	// 构造 client
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("HUNYUAN_API_KEY")),                 // 混元 APIKey
		option.WithBaseURL("https://api.hunyuan.cloud.tencent.com/v1/"), // 混元 endpoint
	)

	appPairStrs, err := a.picStr2Urls(param.Pictures)
	if err != nil {
		return false, err
	}

	urls := appPairStrs.Urls
	for _, url := range urls {
		// 自定义参数传参示例
		completion, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.ChatCompletionUserMessageParam{
					Role: openai.F(openai.ChatCompletionUserMessageParamRoleUser),
					Content: openai.F([]openai.ChatCompletionContentPartUnionParam{
						openai.TextPart("这是哪个景区，仅告诉我景区名字，不要过多分析,不要介绍只需要景区名称"),
						openai.ImagePart(url),
					}),
				},
			}),
			Model: openai.F("hunyuan-vision"),
		},
		)

		if err != nil {
			log.Println(err)
		}

		log.Println(completion.Choices[0].Message.Content)

		now := completion.Choices[0].Message.Content
		if now == name {
			return true, nil
		} else {
		}
	}
	return false, nil
}
func (a *AppService) checkValid(ctx context.Context, param entity.Application) entity.Status {
	status := entity.StatusRefused
	statusBool, err := a.reqOpenai(ctx, &param, param.Name)
	if err != nil {
		log.Println(err)
		return entity.StatusRefused
	}
	if statusBool {
		status = entity.StatusAccepted
	} else {
		status = entity.StatusRefused
		return status
	}

	statusBool, err = a.reqPicValid(ctx, &param)
	if err != nil {
		log.Println(err)
		return entity.StatusRefused
	}
	if statusBool {
		status = entity.StatusAccepted
	} else {
		status = entity.StatusRefused
		return status
	}

	statusBool = sceneExtract(param.Introduction, param.Name)
	if statusBool {
		status = entity.StatusAccepted
	} else {
		status = entity.StatusRefused
		return status
	}

	statusBool = scanMGword(param.Introduction)
	if !statusBool {
		status = entity.StatusRefused
		return status
	} else {
		status = entity.StatusAccepted
	}

	return status
}

// Create 创建
func (a *AppService) Create(ctx context.Context, param entity.AppResp, userId int64, status entity.Status) (*entity.AppResp, error) {
	app, err := a.GetAppByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	currentTime := common.NowInLocal()

	pics := param.Pictures
	data, err := json.Marshal(pics)
	if err != nil {
		return nil, err
	}

	hostmanPics := param.HostmanPics
	hostmanData, err := json.Marshal(hostmanPics)
	if err != nil {
		return nil, err
	}

	hotParams := make([]hotBody, 0)
	n := len(param.Hot)

	for i := 0; i < n; i++ {
		hotParams = append(hotParams, hotBody{
			HotName: param.Hot[i].HotName,
			HotPics: param.Hot[i].HotPics,
			HotDesc: param.Hot[i].HotDesc,
		})
	}
	hot, err := json.Marshal(hotParams)
	if err != nil {
		return nil, err
	}

	enBody, err := a.TranslateWord(ctx, param)
	if err != nil {
		return nil, err
	}

	en, err := json.Marshal(enBody)
	if err != nil {
		return nil, err
	}

	params := entity.Application{
		Pictures:     string(data),
		HostmanPics:  string(hostmanData),
		Hot:          string(hot),
		Name:         param.Name,
		Address:      param.Address,
		Describe:     param.Describe,
		Phone:        param.Phone,
		Rate:         param.Rate,
		Introduction: param.Introduction,
		PostCode:     param.PostCode,
		PostName:     param.PostName,
		WorkBeginAt:  param.WorkBeginAt,
		WorkEndAt:    param.WorkEndAt,
		HaveVege:     param.HaveVege,
		HostmanName:  param.HostmanName,
		HostmanThink: param.HostmanThink,
		Yuyue:        param.Yuyue,
		Leixing:      param.Leixing,
		Per:          param.Per,
		En:           string(en),
	}

	dates, err := json.Marshal(param.WorkDate)
	if err != nil {
		return nil, err
	}

	params.WorkDate = string(dates)

	log.Printf("params = %v\n", params)

	if app != nil {
		params.Id = app.Id
		params.UserId = app.UserId
		status = a.checkValid(ctx, params)
		params.Status = status
		params.UpdatedAt = currentTime
		params.CreatedAt = app.CreatedAt
		tempApp, err := a.Update(ctx, userId, app.Id, &params)
		if err != nil {
			return nil, err
		}

		ret, err := a.application2Resp(tempApp)
		if err != nil {
			return nil, err
		}
		return ret, nil
	}

	params.UserId = userId
	params.CreatedAt = currentTime
	params.UpdatedAt = currentTime
	status = a.checkValid(ctx, params)
	params.Status = status

	tempApp, err := a.db.Create(ctx, &params)
	if err != nil {
		return nil, err
	}

	ret, err := a.application2Resp(tempApp)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// List 列表
func (a *AppService) List(ctx context.Context, offset, limit int, lang string) (map[string]interface{}, error) {
	apps, total, err := a.db.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	params, err := a.apps2Resp(apps)
	if err != nil {
		return nil, err
	}

	if lang == "en" {
		return map[string]interface{}{
			"total": total,
			"apps":  a.multiZh2En(params),
		}, nil
	}

	return map[string]interface{}{
		"total": total,
		"apps":  params,
	}, nil
}

// List 列表
func (a *AppService) SearchList(ctx context.Context, offset, limit int) (map[string]interface{}, error) {
	apps, total, err := a.db.SearchList(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	params, err := a.apps2Resp(apps)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total": total,
		"apps":  params,
	}, nil
}

func (as *AppService) GetAppByUserId(ctx context.Context, userId int64) (*entity.AppResp, error) {
	app, err := as.db.GetByUserId(ctx, userId)
	if err != nil {
		log.Printf("GetAppByUserId_err: %+v", err)
		return nil, err
	}

	if app == nil {
		return nil, nil
	}

	param, err := as.application2Resp(app)
	if err != nil {
		return nil, err
	}

	return param, nil
}

func (as *AppService) ReadApp(ctx context.Context, id int64) (*entity.AppResp, error) {
	app, err := as.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if app.Status != entity.StatusAccepted {
		return nil, pkgs.ServerError
	}

	param, err := as.application2Resp(app)
	if err != nil {
		return nil, err
	}

	return param, nil
}

// Get 获取自己的应用
func (a *AppService) GetSelfApp(ctx context.Context, userId, id int64) (*entity.AppResp, error) {
	app, err := a.db.Get(ctx, id)
	if err != nil {
		log.Printf("Queryapp_Get_err: %+v", err)
		return nil, fmt.Errorf("app not found")
	}
	if app == nil {
		return nil, pkgs.AppNotFound
	}

	user, err := a.account.GetAccount(ctx, userId)
	if err != nil {
		return nil, err
	}

	if app.UserId != userId && user.Status != entity.AdministratorStatus {
		return nil, pkgs.NoPermission
	}

	param, err := a.application2Resp(app)
	if err != nil {
		return nil, err
	}

	return param, nil
}

func (as *AppService) GetAppName(ctx context.Context, userId, id int64) (string, error) {
	app, err := as.GetApp(ctx, id)
	if err != nil {
		return "", err
	}

	if app.UserId != userId {
		return "", pkgs.NoPermission
	}

	return app.Name, nil
}

func (as *AppService) GetAppDescribe(ctx context.Context, userId, id int64) (string, error) {
	app, err := as.GetApp(ctx, id)
	if err != nil {
		return "", err
	}

	if app.UserId != userId {
		return "", pkgs.NoPermission
	}

	return app.Describe, nil
}

func (as *AppService) GetAppAddress(ctx context.Context, userId, id int64) (string, error) {
	app, err := as.GetApp(ctx, id)
	if err != nil {
		return "", err
	}

	if app.UserId != userId {
		return "", pkgs.NoPermission
	}

	return app.Address, nil
}

func (as *AppService) GetAppPhone(ctx context.Context, userId, id int64) (string, error) {
	app, err := as.GetApp(ctx, id)
	if err != nil {
		return "", err
	}

	if app.UserId != userId {
		return "", pkgs.NoPermission
	}

	return app.Phone, nil
}

func (a *AppService) Delete(ctx context.Context, userId, id int64) error {
	app, err := a.db.Get(ctx, id)
	if err != nil {
		return err
	}

	if app.UserId != userId {
		return pkgs.NoPermission
	}

	return a.db.Delete(ctx, id)
}

func (a *AppService) Update(ctx context.Context, userId, id int64, param *entity.Application) (*entity.Application, error) {
	app, err := a.GetApp(ctx, id)
	if err != nil {
		return nil, err
	}

	if app.UserId != userId {
		return nil, pkgs.NoPermission
	}
	return a.db.Update(ctx, param)
}

func (a *AppService) UpdateStatus(ctx context.Context, userId, id int64, status string) error {
	app, err := a.GetApp(ctx, id)
	if err != nil {
		return err
	}

	user, err := a.account.GetAccount(ctx, userId)
	if err != nil {
		return err
	}
	accountStatus := user.Status
	if accountStatus != entity.AdministratorStatus {
		return pkgs.NoPermission
	}

	if status == "通过" {
		app.Status = entity.StatusAccepted
	} else if status == "拒绝" {
		app.Status = entity.StatusRefused
	}

	param, err := a.appResp2Application(app)
	if err != nil {
		return err
	}
	_, err = a.db.Update(ctx, param)

	return err
}

func (as *AppService) SearchByName(ctx context.Context, offset, limit int, name string) (map[string]interface{}, error) {
	apps, total, err := as.db.SearchByName(ctx, offset, limit, name)
	if err != nil {
		return nil, err
	}

	params, err := as.apps2Resp(apps)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total": total,
		"apps":  params,
	}, nil
}

func (as *AppService) SearchByCityAndName(ctx context.Context, offset, limit int, name, city string) (map[string]interface{}, error) {
	apps, total, err := as.db.SearchByCityAndName(ctx, offset, limit, name, city)
	if err != nil {
		return nil, err
	}

	params, err := as.apps2Resp(apps)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total": total,
		"apps":  params,
	}, nil
}

func (a *AppService) apps2Resp(apps []entity.Application) ([]entity.AppResp, error) {
	params := make([]entity.AppResp, 0)
	for _, app := range apps {
		param := &entity.AppResp{
			Id:           app.Id,
			UserId:       app.UserId,
			Name:         app.Name,
			Address:      app.Address,
			Describe:     app.Describe,
			Rate:         app.Rate,
			Phone:        app.Phone,
			Introduction: app.Introduction,
			PostCode:     app.PostCode,
			PostName:     app.PostName,
			WorkBeginAt:  app.WorkBeginAt,
			WorkEndAt:    app.WorkEndAt,
			HaveVege:     app.HaveVege,
			Status:       app.Status,
			CreatedAt:    app.CreatedAt,
			UpdatedAt:    app.UpdatedAt,
			HostmanName:  app.HostmanName,
			HostmanThink: app.HostmanThink,
			Yuyue:        app.Yuyue,
			Leixing:      app.Leixing,
			Per:          app.Per,
		}

		var dates []string
		err := json.Unmarshal([]byte(app.WorkDate), &dates)
		if err != nil {
			return nil, err
		}

		param.WorkDate = dates

		appPairStrs, err := a.picStr2Urls(app.Pictures)
		if err != nil {
			return nil, err
		}

		param.Pictures = appPairStrs.Pics
		param.PicsUrl = appPairStrs.Urls

		var hot []hotBody
		err = json.Unmarshal([]byte(app.Hot), &hot)
		if err != nil {
			return nil, err
		}

		hotParam := make([]entity.HotBody, 0)
		for _, hotbody := range hot {
			hotPictures, err := json.Marshal(hotbody.HotPics)
			if err != nil {
				return nil, err
			}
			hotPairStrs, err := a.picStr2Urls(string(hotPictures))
			if err != nil {
				return nil, err
			}

			hotParam = append(hotParam, entity.HotBody{
				HotName:    hotbody.HotName,
				HotDesc:    hotbody.HotDesc,
				HotPics:    hotPairStrs.Pics,
				HotPicUrls: hotPairStrs.Urls,
			})
		}

		param.Hot = hotParam

		hostmanPictures, err := a.picStr2Urls(app.HostmanPics)
		if err != nil {
			return nil, err
		}

		param.HostmanPics = hostmanPictures.Pics
		param.HostmanUrls = hostmanPictures.Urls

		var enBody entity.AppEn
		err = json.Unmarshal([]byte(app.En), &enBody)
		if err != nil {
			return nil, err
		}

		param.En = enBody

		params = append(params, *param)
	}

	return params, nil
}

// List 列表
func (a *AppService) SearchByCity(ctx context.Context, offset, limit int, key string) (map[string]interface{}, error) {
	apps, total, err := a.db.SearchByCity(ctx, offset, limit, key)
	if err != nil {
		return nil, err
	}

	params, err := a.apps2Resp(apps)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total": total,
		"apps":  params,
	}, nil
}

func (as *AppService) DownloadDoc(ctx context.Context, userId, id int64) (string, error) {
	user, err := as.account.GetAccount(ctx, userId)
	if err != nil {
		return "", err
	}

	if user.Status != entity.AdministratorStatus {
		return "", pkgs.NoPermission
	}

	app, err := as.db.Get(ctx, id)
	if err != nil {
		return "", err
	}

	params, err := as.application2Resp(app)
	if err != nil {
		return "", fmt.Errorf("download_doc_application_to_response_err: %v", err)
	}

	doc := document.New()
	defer doc.Close()

	para := doc.AddParagraph()
	run := para.AddRun()
	run.AddText(fmt.Sprintf("景点名称:\t%s", params.Name))
	run.AddBreak()

	// run.AddText(fmt.Sprintf("餐厅类型:\t%s", params.Leixing))
	// run.AddBreak()

	run.AddText(fmt.Sprintf("门票价:\t%s", params.Per))
	run.AddBreak()

	run.AddText(fmt.Sprintf("详细地址:\t%s", params.Address))
	run.AddBreak()

	run.AddText(fmt.Sprintf("联系方式:\t%s", params.Phone))
	run.AddBreak()

	run.AddText(fmt.Sprintf("营业日期:\t%s", params.WorkDate))
	run.AddBreak()

	run.AddText(fmt.Sprintf("营业时间:\t%s~%s", params.WorkBeginAt, params.WorkEndAt))
	run.AddBreak()

	// run.AddText(fmt.Sprintln("是否可以预约:\t"))
	// if params.Yuyue == 1 {
	// 	run.AddText("可以")
	// } else {
	// 	run.AddText("否")
	// }
	// run.AddBreak()

	// run.AddText(fmt.Sprintln("是否有素食:\t"))
	// if params.HaveVege == 1 {
	// 	run.AddText("是")
	// } else {
	// 	run.AddText("否")
	// }

	run.AddText("置顶图片:")
	run.AddBreak()

	tempDir := "../tempImages"
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	imgs := []officedocCommon.ImageRef{}
	for i, url := range params.PicsUrl {
		response, err := http.Get(url)
		if err != nil {
			log.Fatalf("Error downloading image: %v", err)
		}
		defer response.Body.Close()

		imgPath := filepath.Join(tempDir, filepath.Base(params.Pictures[i]))
		file, err := os.Create(imgPath)
		if err != nil {
			log.Fatalf("Error creating image file: %v", err)
		}
		defer file.Close()

		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatalf("Error saving image file: %v", err)
		}

		img, err := officedocCommon.ImageFromFile(imgPath)
		if err != nil {
			return "", fmt.Errorf("unable to create image: %s", err)
		}

		imgref, err := doc.AddImage(img)
		if err != nil {
			return "", err
		}
		imgs = append(imgs, imgref)
	}

	for _, img := range imgs {
		inl, err := run.AddDrawingInline(img)
		if err != nil {
			return "", fmt.Errorf("unable to add inline image: %s", err)
		}
		inl.SetSize(2*measurement.Inch, 2*measurement.Inch)
		run.AddTab()
	}
	run.AddBreak()

	run.AddText(fmt.Sprintf("简介:\t%s", params.Introduction))
	// run.AddText(fmt.Sprintf("详细介绍:\t%s", params.Describe))

	// tempHostDir := "../tempHostImages"
	// os.MkdirAll(tempHostDir, os.ModePerm)
	// defer os.RemoveAll(tempHostDir)

	// hostImgs := []officedocCommon.ImageRef{}
	// for i, url := range params.HostmanUrls {
	// 	response, err := http.Get(url)
	// 	if err != nil {
	// 		log.Fatalf("Error downloading image: %v", err)
	// 	}
	// 	defer response.Body.Close()

	// 	imgPath := filepath.Join(tempHostDir, filepath.Base(params.HostmanPics[i]))
	// 	file, err := os.Create(imgPath)
	// 	if err != nil {
	// 		log.Fatalf("Error creating image file: %v", err)
	// 	}
	// 	defer file.Close()

	// 	_, err = io.Copy(file, response.Body)
	// 	if err != nil {
	// 		log.Fatalf("Error saving image file: %v", err)
	// 	}

	// 	img, err := officedocCommon.ImageFromFile(imgPath)
	// 	if err != nil {
	// 		return "", fmt.Errorf("unable to create image: %s", err)
	// 	}

	// 	imgref, err := doc.AddImage(img)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	hostImgs = append(hostImgs, imgref)
	// }

	// run.AddText(fmt.Sprintf("主理人名字:\t%s", params.HostmanName))
	// run.AddBreak()

	// for _, img := range hostImgs {
	// 	inl, err := run.AddDrawingInline(img)
	// 	if err != nil {
	// 		return "", fmt.Errorf("unable to add inline image: %s", err)
	// 	}
	// 	inl.SetSize(2*measurement.Inch, 2*measurement.Inch)
	// 	run.AddTab()
	// }
	// run.AddBreak()

	// run.AddText(fmt.Sprintf("主理人理念:\t%s", params.HostmanThink))
	// run.AddBreak()

	// hotLen := len(params.Hot)

	// for j := 0; j < hotLen; j++ {
	// 	run.AddText(fmt.Sprintf("招牌菜 %d: \t%s", j, params.Hot[j].HotName))
	// 	run.AddBreak()

	// 	tempHotDir := fmt.Sprintf("../tempHotImages_%d", j+1)
	// 	os.MkdirAll(tempHotDir, os.ModePerm)
	// 	defer os.RemoveAll(tempHotDir)
	// 	hotImgs := []officedocCommon.ImageRef{}
	// 	for i, url := range params.Hot[j].HotPicUrls {
	// 		response, err := http.Get(url)
	// 		if err != nil {
	// 			log.Fatalf("Error downloading image: %v", err)
	// 		}
	// 		defer response.Body.Close()

	// 		imgPath := filepath.Join(tempHotDir, filepath.Base(params.Hot[j].HotPics[i]))
	// 		file, err := os.Create(imgPath)
	// 		if err != nil {
	// 			log.Fatalf("Error creating image file: %v", err)
	// 		}
	// 		defer file.Close()

	// 		_, err = io.Copy(file, response.Body)
	// 		if err != nil {
	// 			log.Fatalf("Error saving image file: %v", err)
	// 		}

	// 		img, err := officedocCommon.ImageFromFile(imgPath)
	// 		if err != nil {
	// 			return "", fmt.Errorf("unable to create image: %s", err)
	// 		}

	// 		imgref, err := doc.AddImage(img)
	// 		if err != nil {
	// 			return "", err
	// 		}
	// 		hotImgs = append(hotImgs, imgref)
	// 	}

	// 	for _, img := range hotImgs {
	// 		inl, err := run.AddDrawingInline(img)
	// 		if err != nil {
	// 			return "", fmt.Errorf("unable to add inline image: %s", err)
	// 		}
	// 		inl.SetSize(2*measurement.Inch, 2*measurement.Inch)
	// 		run.AddTab()
	// 	}
	// 	run.AddBreak()

	// 	run.AddText(fmt.Sprintf("招牌菜描述:\t%s", params.Hot[j].HotDesc))
	// 	run.AddBreak()
	// }

	saveFilePath := fmt.Sprintf("%s.docx", params.Name)
	err = doc.SaveToFile(saveFilePath)
	defer os.Remove(saveFilePath)
	if err != nil {
		return "", err
	}

	err = pkgs.UploadOSSByFilePath(saveFilePath, saveFilePath, false)
	if err != nil {
		return "", fmt.Errorf("Upload_OSS_File_err:%v", err)
	}

	signedUrl, _ := pkgs.SignedUrl(saveFilePath, false)
	return signedUrl, nil
}
