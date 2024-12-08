package application

import (
	"context"
	"restaurant/domain"
	"restaurant/entity"
	"restaurant/pkgs"
)

type StaticBannerPicService struct {
	service *domain.StaticBannerPicRepo
	as      AccountService
}

func NewStaticBannerPicService() *StaticBannerPicService {
	return &StaticBannerPicService{
		service: domain.NewStaticBannerPicRepo(),
		as:      *NewAccountService(),
	}
}

func (ss *StaticBannerPicService) Get(ctx context.Context, id int64) (map[string]interface{}, error) {
	app, err := ss.service.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	picture := app.Picture
	picUrl, _ := pkgs.SignedUrl(picture, false)
	return map[string]interface{}{
		"id":          app.Id,
		"pic_file": 	 app.Picture,
		"picture":     picUrl,
		"redirct_url": app.RedirectUrl,
	}, nil
}

func (ss *StaticBannerPicService) List(ctx context.Context) (map[string]interface{}, error) {
	apps, total, err := ss.service.List(ctx)
	if err != nil {
		return nil, err
	}

	type staticBannerPicList struct {
		Id          int64  `json:"id"`
		PicFile     string `json:"pic_file"`
		PicUrl string `json:"pic_url"`
		RedirectUrl string `json:"redirect_url"`
	}
	pics := make([]staticBannerPicList, 0)
	for _, app := range apps {
		picture := app.Picture
		picUrl, _ := pkgs.SignedUrl(picture, false)
		pics = append(pics, staticBannerPicList{
			Id:          app.Id,
			PicFile: picture,
			PicUrl: picUrl,
			RedirectUrl: app.RedirectUrl,
		})
	}
	return map[string]interface{}{
		"total": total,
		"apps":  pics,
	}, nil
}

func (ss *StaticBannerPicService) Update(ctx context.Context, userId, id int64, filename, url string) error {
	user, err := ss.as.GetAccount(ctx, userId)
	if err != nil {
		return err
	}

	if user.Status != entity.AdministratorStatus {
		return pkgs.NoPermission
	}

	return ss.service.Update(ctx, id, map[string]interface{}{
		"picture":      filename,
		"redirect_url": url,
	})
}
