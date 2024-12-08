package application

import (
	"context"
	"log"
	"restaurant/domain"
	"restaurant/entity"
	"restaurant/pkgs"
)

type ItemConfigService struct {
	service *domain.ItemConfigRepo
	as      AccountService
}

func NewItemConfigService() *ItemConfigService {
	return &ItemConfigService{
		service: domain.NewItemConfigRepo(),
		as:      *NewAccountService(),
	}
}

func (is *ItemConfigService) Get(ctx context.Context) (*entity.ItemConfig, error) {
	return is.service.Get(ctx)
}

func (is *ItemConfigService) Update(ctx context.Context, userId int64, status entity.ItemStatus) error {
	user, err := is.as.GetAccount(ctx, userId)
	if err != nil {
		return err
	}

	if user.Status != entity.AdministratorStatus {
		return pkgs.NoPermission
	}

	log.Println(status)

	return is.service.Update(ctx, map[string]interface{}{
		"prize": status,
	})
}

func (is *ItemConfigService) Rename(ctx context.Context, userId int64, param entity.ItemConfig) error {
	user, err := is.as.GetAccount(ctx, userId)
	if err != nil {
		return err
	}

	if user.Status != entity.AdministratorStatus {
		return pkgs.NoPermission
	}

	configMap := make(map[string]interface{})
	if param.AppName != "" {
		configMap["app_name"] = param.AppName
	}
	if param.ChefName != "" {
		configMap["chef_name"] = param.ChefName
	}
	if param.GuestName != "" {
		configMap["guest_name"] = param.GuestName
	}
	if param.NewsName != "" {
		configMap["news_name"] = param.NewsName
	}
	if param.PrizeName != "" {
		configMap["prize_name"] = param.PrizeName
	}
	return is.service.Update(ctx, configMap)
}
