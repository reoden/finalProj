package application

import (
	"context"
	"restaurant/domain"
	"restaurant/entity"
	"restaurant/pkgs"
)

type LogoService struct {
	service *domain.LogoRepo
	as      AccountService
}

func NewLogoService() *LogoService {
	return &LogoService{
		service: domain.NewLogoRepo(),
		as:      *NewAccountService(),
	}
}

func (ls *LogoService) Get(ctx context.Context, id int64) (map[string]interface{}, error) {
	param, err := ls.service.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	pic := param.Pic
	picUrl, err := pkgs.SignedUrl(pic, false)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":param.Id,
		"pic_url": picUrl,
		"pic_file": param.Pic,
	},nil
}

func (ls *LogoService) Update(ctx context.Context, userId int64, param *entity.Logo) error {
	user, err := ls.as.GetAccount(ctx, userId)
	if err != nil {
		return err
	}

	if user.Status != entity.AdministratorStatus {
		return pkgs.NoPermission
	}

	return ls.service.Update(ctx, param)
}

func (ls *LogoService) List(ctx context.Context) (map[string]interface{}, error) {
	apps, total, err := ls.service.List(ctx)
	if err != nil {
		return nil, err
	}

	type logoListBody struct {
		Id int64 `json:"id"`
		PicUrl string `json:"pic_url"`
		PicFile string `json:"pic_file"`
	}
	pics := make([]logoListBody, 0)
	for _, app := range apps {
		picture := app.Pic
		picUrl, _ := pkgs.SignedUrl(picture, false)
		pics = append(pics, logoListBody{
			Id:  app.Id,
			PicUrl: picUrl,
			PicFile: picture,
		})
	}

	return map[string]interface{}{
		"total": total,
		"apps":  pics,
	}, nil
}
