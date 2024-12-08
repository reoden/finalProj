package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"restaurant/common"
	"restaurant/domain"
	"restaurant/entity"
	"restaurant/pkgs"

	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

type ChefService struct {
	db      *domain.ChefRepo
	account AccountService
	chat    ChatService
}

func NewChefService() *ChefService {
	return &ChefService{
		db:      domain.NewChefRepo(),
		account: *NewAccountService(),
		chat:    *NewChatService(),
	}
}

func (as *ChefService) GetApp(ctx context.Context, id int64) (*entity.Chef, error) {
	app, err := as.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (as *ChefService) translate(ctx context.Context, param *entity.ChefResp) (*entity.EnChef, error) {
	nameResult, err := as.chat.Chat(ctx, as.chat.GetParams(param.Name), nil)
	if err != nil {
		return nil, err
	}

	addressResult, err := as.chat.Chat(ctx, as.chat.GetParams(param.Address), nil)
	if err != nil {
		return nil, err
	}

	describeResult, err := as.chat.Chat(ctx, as.chat.GetParams(param.Describe), nil)
	if err != nil {
		return nil, err
	}

	return &entity.EnChef{
		Name:     nameResult.Content,
		Address:  addressResult.Content,
		Describe: describeResult.Content,
	}, nil
}

func (as *ChefService) zh2en(ctx context.Context, param *entity.Chef) (entity.Chef, error) {
	var enBody entity.EnChef
	err := json.Unmarshal([]byte(param.En), &enBody)

	if err != nil {
		return entity.Chef{}, err
	}
	return entity.Chef{
		Name:     enBody.Name,
		Describe: enBody.Describe,
		Address:  enBody.Address,
	}, nil
}

func (as *ChefService) multiZh2En(ctx context.Context, params []entity.Chef) (apps []entity.Chef, err error) {
	for _, param := range params {
		app, err := as.zh2en(ctx, &param)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	return apps, nil
}

func (as *ChefService) ReadApp(ctx context.Context, id int64, lang string) (*entity.ChefResp, error) {
	app, err := as.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if app.Status != entity.StatusAccepted {
		return nil, pkgs.ServerError
	}

	param, err := as.chef2Resp(app)
	if err != nil {
		return nil, err
	}

	return param, nil
}

func (as *ChefService) GetApps(ctx context.Context, ids []int64, lang string) ([]*entity.ChefResp, error) {
	apps, err := as.db.GetApps(ctx, ids)
	if err != nil {
		return nil, err
	}

	params := make([]*entity.ChefResp, 0)
	for _, app := range apps {
		param, err := as.chef2Resp(app)
		if err != nil {
			return nil, err
		}

		params = append(params, param)
	}

	return params, nil
}

func (a *ChefService) chef2Resp(app *entity.Chef) (*entity.ChefResp, error) {
	param := &entity.ChefResp{
		Id:        app.Id,
		Name:      app.Name,
		Address:   app.Address,
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

	var enChef entity.EnChef
	err = json.Unmarshal([]byte(app.En), &enChef)
	if err != nil {
		return nil, err
	}

	param.En = enChef

	return param, nil
}

func (a *ChefService) appResp2Chef(app *entity.ChefResp) (*entity.Chef, error) {
	param := &entity.Chef{
		Id:        app.Id,
		Name:      app.Name,
		Address:   app.Address,
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
func (a *ChefService) Create(ctx context.Context, param entity.ChefResp, userId int64) (*entity.ChefResp, error) {
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

	enChef, err := a.translate(ctx, &param)
	if err != nil {
		return nil, err
	}

	enBody, err := json.Marshal(enChef)
	if err != nil {
		return nil, err
	}

	if app != nil {
		newApp, err := a.appResp2Chef(&param)
		if err != nil {
			return nil, err
		}

		newApp.En = string(enBody)

		newApp.Status = entity.StatusAccepted
		newApp.CreatedAt = app.CreatedAt

		tempApp, err := a.Update(ctx, newApp)
		if err != nil {
			return nil, err
		}

		return a.chef2Resp(tempApp)
	}

	pics := param.Pictures
	data, err := json.Marshal(pics)
	if err != nil {
		return nil, err
	}

	params := entity.Chef{
		Pictures: string(data),
		Name:     param.Name,
		Address:  param.Address,
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

	ret, err := a.chef2Resp(tempApp)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// convert BIG5 to UTF-8
func DecodeBig5(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// convert UTF-8 to BIG5
func EncodeBig5(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewEncoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func pythonExecute(args string) ([]string, error) {
	pycode := fmt.Sprintf("from chinese_stroke_sorting import sort_by_stroke; import json; name_list = %s; result = sort_by_stroke(name_list); print(json.dumps(result))", args)

	log.Println(pycode)
	arg := []string{
		"-c",
		pycode,
	}

	cmd := exec.Command("python3", arg...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Error when executing python: %s\n", stderr.Bytes())
		return nil, fmt.Errorf("error executing python: %w", err)
	}

	res := []string{}
	err = json.Unmarshal(stdout.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func trans(s string) string {
	dic := make(map[rune]rune)
	dic['0'] = '〇'
	dic['1'] = '一'
	dic['2'] = '二'
	dic['3'] = '三'
	dic['4'] = '四'
	dic['5'] = '五'
	dic['6'] = '六'
	dic['7'] = '七'
	dic['8'] = '八'
	dic['9'] = '九'

	ret := ""
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			ret = ret + string(dic[ch])
		} else if ch >= 'a' && ch <= 'z' {
			cur := int(ch - 'a')
			a, b := cur/10, cur%10
			ret = ret + string(dic[rune(a+'0')]) + string(dic[rune(b+'0')])
		} else if ch >= 'A' && ch <= 'Z' {
			cur := int(ch - 'A')
			a, b := cur/10, cur%10
			ret = ret + string(dic[rune(a+'0')]) + string(dic[rune(b+'0')])
		} else {
			ret = ret + string(ch)
		}
	}

	return ret
}

func (gs *ChefService) sortByStroke(params []*entity.ChefResp) ([]int64, error) {
	names := make([]string, 0)
	for _, param := range params {
		name := param.Name
		name = trans(name)
		names = append(names, name)
	}

	log.Println(names)
	index := make(map[string]int64)
	for i := range names {
		index[names[i]] = int64(i)
	}

	args, err := json.Marshal(names)

	if err != nil {
		return nil, err
	}
	names, err = pythonExecute(string(args))
	if err != nil {
		return nil, err
	}

	ids := []int64{}
	for _, name := range names {
		ids = append(ids, index[name])
	}

	return ids, nil
}

// List 列表
func (a *ChefService) List(ctx context.Context, lang string) (map[string]interface{}, error) {
	apps, total, err := a.db.List(ctx)
	if err != nil {
		return nil, err
	}

	params := make([]*entity.ChefResp, 0)
	for _, app := range apps {
		param, err := a.chef2Resp(&app)
		if err != nil {
			return nil, err
		}

		params = append(params, param)
	}

	ids, err := a.sortByStroke(params)

	if err != nil {
		return nil, err
	}
	cur := make([]*entity.ChefResp, 0)

	for _, index := range ids {
		cur = append(cur, params[index])
	}

	return map[string]interface{}{
		"total": total,
		"apps":  cur,
	}, nil
}

func (as *ChefService) Update(ctx context.Context, param *entity.Chef) (*entity.Chef, error) {
	return as.db.Update(ctx, param)
}

func (a *ChefService) Delete(ctx context.Context, userId, id int64) error {
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

func (a *ChefService) Get(ctx context.Context, userId, id int64, lang string) (*entity.ChefResp, error) {
	user, err := a.account.GetAccount(ctx, userId)
	if err != nil {
		return nil, err
	}

	if user.Status != entity.AdministratorStatus {
		return nil, pkgs.NoPermission
	}

	chef, err := a.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	param, err := a.chef2Resp(chef)
	if err != nil {
		return nil, err
	}

	return param, nil
}
