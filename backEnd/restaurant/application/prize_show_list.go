package application

import (
	"context"
	"errors"
	"log"
	"restaurant/common"
	"restaurant/domain"
	"restaurant/entity"
	"restaurant/pkgs"
)

type PrizeShowListService struct {
	db      *domain.PrizeShowListRepo
	as      AppService
	account AccountService
}

func NewPrizeShowListService() *PrizeShowListService {
	return &PrizeShowListService{
		db:      domain.NewPrizeShowListRepo(),
		as:      *NewAppService(),
		account: *NewAccountService(),
	}
}

func (br *PrizeShowListService) Create(ctx context.Context, userId int64, param *entity.BannerBody) (*entity.BannerBody, error) {
	user, err := br.account.GetAccount(ctx, userId)
	if err != nil {
		return nil, err
	}

	if user.Status != entity.AdministratorStatus {
		return nil, pkgs.NoPermission
	}

	app, err := br.as.GetApp(ctx, param.AppId)
	if err != nil {
		return nil, err
	}

	if app.Status != entity.StatusAccepted {
		return nil, errors.New("当前商家尚未通过审核")
	}

	banner, err := br.db.GetByAppId(ctx, param.AppId)
	if err != nil {
		return nil, err
	}

	if banner != nil {
		return nil, errors.New("重复添加")
	}

	currentTime := common.NowInLocal()
	param.CreatedAt = currentTime
	param.Status = entity.BannerShowStatus
	return br.db.Create(ctx, param)
}

func (br *PrizeShowListService) CreateBanners(ctx context.Context, param []*entity.BannerBody) ([]*entity.BannerBody, error) {
	currentTime := common.NowInLocal()
	for i := range param {
		param[i].CreatedAt = currentTime
		param[i].Status = entity.BannerShowStatus
	}
	return br.db.CreateBanners(ctx, param)
}

func (br *PrizeShowListService) Delete(ctx context.Context, userId, id int64) error {
	user, err := br.account.GetAccount(ctx, userId)
	if err != nil {
		return err
	}
	accountStatus := user.Status
	if accountStatus != entity.AdministratorStatus {
		return pkgs.NoPermission
	}

	return br.db.Delete(ctx, id)
}

func (br *PrizeShowListService) DeleteIds(ctx context.Context, ids []int64) error {
	return br.db.DeleteIds(ctx, ids)
}

func (br *PrizeShowListService) List(ctx context.Context) (*entity.BannerListResp, error) {
	apps, _, err := br.db.List(ctx)
	if err != nil {
		return nil, err
	}

	appIds := make([]int64, 0)
	for _, app := range apps {
		appIds = append(appIds, app.AppId)
	}

	log.Println(appIds)

	if len(appIds) == 0 {
		return nil, nil
	}

	params, err := br.as.GetAppsOrder(ctx, appIds)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0)
	removeIds := make([]int64, 0)
	for _, param := range params {
		if param.Status == entity.StatusAccepted {
			ids = append(ids, param.Id)
		} else {
			removeIds = append(removeIds, param.Id)
		}
	}

	if len(removeIds) != 0 {
		err = br.DeleteIds(ctx, removeIds)
		if err != nil {
			return nil, err
		}
	}

	bannerResp := entity.BannerListResp{}
	if len(ids) == 0 {
		return &bannerResp, nil
	}

	params, err = br.as.GetAppsOrder(ctx, ids)
	if err != nil {
		return nil, err
	}

	bannerResp = entity.BannerListResp{
		Apps:  params,
		Total: int64(len(params)),
	}

	return &bannerResp, nil
}

func (br *PrizeShowListService) Update(ctx context.Context, userId int64, appIds []int64) error {
	user, err := br.account.GetAccount(ctx, userId)
	if err != nil {
		return err
	}
	if user.Status != entity.AdministratorStatus {
		return pkgs.NoPermission
	}

	banners, err := br.db.GetAll(ctx)
	if err != nil {
		return err
	}

	bannerIds := make([]int64, 0, len(banners))
	for _, banner := range banners {
		bannerIds = append(bannerIds, banner.Id)
	}

	bannerBodys := make([]*entity.BannerBody, 0)
	for _, appId := range appIds {
		bannerBodys = append(bannerBodys, &entity.BannerBody{
			AppId: appId,
		})
	}

	_, err = br.CreateBanners(ctx, bannerBodys)
	if err != nil {
		return err
	}

	return br.db.DeleteBanners(ctx, bannerIds)
}

func (br *PrizeShowListService) SearchByCity(ctx context.Context, key string) (map[string]interface{}, error) {
	apps, total, err := br.db.SearchByCity(ctx, key)
	if err != nil {
		return nil, err
	}

	params, err := br.as.apps2Resp(apps)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total": total,
		"apps":  params,
	}, nil
}

func (br *PrizeShowListService) SearchByName(ctx context.Context, name string) (map[string]interface{}, error) {
	apps, total, err := br.db.SearchByName(ctx, name)
	if err != nil {
		return nil, err
	}

	params, err := br.as.apps2Resp(apps)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total": total,
		"apps":  params,
	}, nil
}

func (br *PrizeShowListService) SearchByCityAndName(ctx context.Context, name, city string) (map[string]interface{}, error) {
	apps, total, err := br.db.SearchByCityAndName(ctx, name, city)
	if err != nil {
		return nil, err
	}

	params, err := br.as.apps2Resp(apps)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total": total,
		"apps":  params,
	}, nil
}

func (br *PrizeShowListService) SearchById(ctx context.Context, id int64) (map[string]interface{}, error) {
	apps, total, err := br.db.SearchById(ctx, id)
	if err != nil {
		return nil, err
	}

	params, err := br.as.apps2Resp(apps)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total": total,
		"apps":  params,
	}, nil
}

func (br *PrizeShowListService) GetPrizeBannerAdver(ctx context.Context) map[string]string {

	banner, err := br.db.GetBanner(ctx)
	if err != nil {
		log.Println(err)
		return nil
	}

	signedUrl, _ := pkgs.SignedUrl(banner.Picture, false)

	return map[string]string{
		"pciture":      banner.Picture,
		"pic_url":      signedUrl,
		"redirect_url": banner.RedirectUrl,
	}
}

func (br *PrizeShowListService) UpdatePrizeBannerAdver(ctx context.Context, userId int64, param entity.PrizeBanner) error {
	user, err := br.account.GetAccount(ctx, userId)
	if err != nil {
		return err
	}

	if user.Status != entity.AdministratorStatus {
		return pkgs.NoPermission
	}

	log.Println(param)
	br.db.UpdateBanner(ctx, map[string]interface{}{
		"picture":      param.Picture,
		"redirect_url": param.RedirectUrl,
	})
	return nil
}
