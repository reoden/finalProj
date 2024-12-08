package domain

import (
	"context"
	"restaurant/entity"
	"restaurant/pkgs"
	"restaurant/pkgs/sferror"

	"gorm.io/gorm"
)

const TableLogo = "logo"

type LogoRepo struct {
	db *gorm.DB
}

func NewLogoRepo() *LogoRepo {
	return &LogoRepo{
		db: pkgs.GetDB(),
	}
}

func (lr *LogoRepo) Get(ctx context.Context, id int64) (*entity.Logo, error) {
	var param *entity.Logo
	err := lr.db.Table(TableLogo).Where("id = ?", id).First(&param).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return param, nil
}

func (lr *LogoRepo) Update(ctx context.Context, param *entity.Logo) error {
	err := lr.db.Table(TableLogo).Where("id = ?", param.Id).Updates(param).Error
	if err != nil {
		return sferror.WithStack(err)
	}

	return nil
}

func (lr *LogoRepo) List(ctx context.Context) ([]*entity.Logo, int64, error) {
	var app []*entity.Logo
	query := lr.db.Table(TableLogo).Where("pic != ''")

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, sferror.WithStack(err)
	}

	err = query.Order("id desc").Find(&app).Error
	if err == gorm.ErrRecordNotFound {
		return nil, 0, nil
	}

	if err != nil {
		return nil, 0, sferror.WithStack(err)
	}

	return app, total, nil
}
