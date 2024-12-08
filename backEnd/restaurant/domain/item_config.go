package domain

import (
	"context"
	"restaurant/entity"
	"restaurant/pkgs"
	"restaurant/pkgs/sferror"

	"gorm.io/gorm"
)

const (
	TableConfig = "config"
)

type ItemConfigRepo struct {
	db *gorm.DB
}

func NewItemConfigRepo() *ItemConfigRepo {
	return &ItemConfigRepo{
		db: pkgs.GetDB(),
	}
}

func (ir *ItemConfigRepo) Get(ctx context.Context) (*entity.ItemConfig, error) {
	var app entity.ItemConfig
	err := ir.db.Table(TableConfig).Where("id = ?", 1).First(&app).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return &app, nil
}

func (ir *ItemConfigRepo) Update(ctx context.Context, param map[string]interface{}) error {
	err := ir.db.Table(TableConfig).Where("id = ?", 1).Updates(&param).Error
	if err != nil {
		return sferror.WithStack(err)
	}

	return nil
}
