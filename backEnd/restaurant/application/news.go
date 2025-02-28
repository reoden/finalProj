package application

import (
	"context"
	"encoding/json"
	"log"
	"restaurant/common"
	"restaurant/domain"
	"restaurant/entity"
	"restaurant/pkgs"
)

type NewsService struct {
	db      *domain.NewsRepo
	account AccountService
	chat    ChatService
}

func NewNewsService() *NewsService {
	return &NewsService{
		db:      domain.NewNewsRepo(),
		account: *NewAccountService(),
		chat:    *NewChatService(),
	}
}

func (as *NewsService) GetApp(ctx context.Context, id int64) (*entity.News, error) {
	app, err := as.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (as *NewsService) translate(ctx context.Context, param *entity.NewsResp) (*entity.EnNews, error) {
	nameResult, err := as.chat.Chat(ctx, as.chat.GetParams(param.Name), nil)
	if err != nil {
		return nil, err
	}

	describeResult, err := as.chat.Chat(ctx, as.chat.GetParams(param.Describe), nil)
	if err != nil {
		return nil, err
	}

	return &entity.EnNews{
		Name:     nameResult.Content,
		Describe: describeResult.Content,
	}, nil
}

func (as *NewsService) zh2en(ctx context.Context, param *entity.News) (entity.News, error) {
	var enBody entity.EnGuest
	err := json.Unmarshal([]byte(param.En), &enBody)

	if err != nil {
		return entity.News{}, err
	}
	return entity.News{
		Name:     enBody.Name,
		Describe: enBody.Describe,
	}, nil
}

func (as *NewsService) ReadApp(ctx context.Context, id int64, lang string) (*entity.NewsResp, error) {
	app, err := as.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if app.Status != entity.StatusAccepted {
		return nil, pkgs.ServerError
	}

	param, err := as.news2Resp(app)
	if err != nil {
		return nil, err
	}

	return param, nil
}

func (as *NewsService) GetApps(ctx context.Context, ids []int64, lang string) ([]*entity.NewsResp, error) {
	apps, err := as.db.GetApps(ctx, ids)
	if err != nil {
		return nil, err
	}

	params := make([]*entity.NewsResp, 0)
	for _, app := range apps {
		param, err := as.news2Resp(app)
		if err != nil {
			return nil, err
		}

		params = append(params, param)
	}

	return params, nil
}

func (a *NewsService) news2Resp(app *entity.News) (*entity.NewsResp, error) {
	param := &entity.NewsResp{
		Id:        app.Id,
		Name:      app.Name,
		Describe:  app.Describe,
		Status:    app.Status,
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
	}

	bodyPics := make([]string, 0)
	var appPics []string
	err := json.Unmarshal([]byte(app.Pictures), &appPics)
	if err != nil {
		return nil, err
	}

	for _, pic := range appPics {
		signedUrl, _ := pkgs.SignedUrl(pic, false)
		bodyPics = append(bodyPics, signedUrl)
	}

	log.Println(bodyPics)
	param.PicUrl = bodyPics
	param.Pictures = appPics

	var en entity.EnNews
	err = json.Unmarshal([]byte(app.En), &en)
	if err != nil {
		return nil, err
	}

	param.En = en

	return param, nil
}

func (a *NewsService) appResp2News(app *entity.NewsResp) (*entity.News, error) {
	param := &entity.News{
		Id:        app.Id,
		Name:      app.Name,
		Describe:  app.Describe,
		Status:    app.Status,
		CreatedAt: app.CreatedAt,
		UpdatedAt: app.UpdatedAt,
	}

	pics, err := json.Marshal(app.Pictures)
	if err != nil {
		return nil, err
	}

	param.Pictures = string(pics)

	en, err := json.Marshal(app.En)
	if err != nil {
		return nil, err
	}

	param.En = string(en)

	return param, nil
}

// Create 创建
func (a *NewsService) Create(ctx context.Context, param entity.NewsResp, userId int64) (*entity.NewsResp, error) {
	user, err := a.account.GetAccount(ctx, userId)
	if err != nil {
		return nil, err
	}

	if user.Status != entity.AdministratorStatus {
		return nil, pkgs.NoPermission
	}

	currentTime := common.NowInLocal()

	app, err := a.db.Get(ctx, param.Id)
	if err != nil {
		return nil, err
	}

	enNews, err := a.translate(ctx, &param)
	if err != nil {
		return nil, err
	}

	enBody, err := json.Marshal(enNews)
	if err != nil {
		return nil, err
	}

	if app != nil {
		newApp, err := a.appResp2News(&param)
		if err != nil {
			return nil, err
		}

		newApp.Status = entity.StatusAccepted
		newApp.En = string(enBody)
		newApp.CreatedAt = app.CreatedAt

		tempApp, err := a.Update(ctx, newApp)
		if err != nil {
			return nil, err
		}

		return a.news2Resp(tempApp)
	}

	pics := param.Pictures
	data, err := json.Marshal(pics)
	if err != nil {
		return nil, err
	}

	params := entity.News{
		Pictures: string(data),
		Name:     param.Name,
		Describe: param.Describe,
		En:       string(enBody),
		Status:   entity.StatusAccepted,
	}

	params.CreatedAt = currentTime
	params.UpdatedAt = currentTime

	tempApp, err := a.db.Create(ctx, &params)
	if err != nil {
		return nil, err
	}

	ret, err := a.news2Resp(tempApp)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// List 列表
func (a *NewsService) List(ctx context.Context, lang string) (map[string]interface{}, error) {
	apps, total, err := a.db.List(ctx)
	if err != nil {
		return nil, err
	}

	params := make([]*entity.NewsResp, 0)
	for _, app := range apps {
		param, err := a.news2Resp(&app)
		if err != nil {
			return nil, err
		}

		params = append(params, param)
	}

	return map[string]interface{}{
		"total": total,
		"apps":  params,
	}, nil
}

func (as *NewsService) Update(ctx context.Context, param *entity.News) (*entity.News, error) {
	return as.db.Update(ctx, param)
}

func (a *NewsService) Delete(ctx context.Context, userId, id int64) error {
	user, err := a.account.GetAccount(ctx, userId)
	if err != nil {
		return err
	}

	if user.Status != entity.AdministratorStatus {
		return pkgs.NoPermission
	}

	_, err = a.db.Get(ctx, id)
	if err != nil {
		return err
	}

	return a.db.Delete(ctx, id)
}

func (a *NewsService) Get(ctx context.Context, userId, id int64, lang string) (*entity.NewsResp, error) {
	user, err := a.account.GetAccount(ctx, userId)
	if err != nil {
		return nil, err
	}

	if user.Status != entity.AdministratorStatus {
		return nil, pkgs.NoPermission
	}

	news, err := a.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	param, err := a.news2Resp(news)
	if err != nil {
		return nil, err
	}

	return param, nil
}
