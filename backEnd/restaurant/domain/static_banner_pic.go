package domain

import (
	"context"
	"restaurant/entity"
	"restaurant/pkgs"
	"restaurant/pkgs/sferror"

	"gorm.io/gorm"
)

type StaticBannerPicRepo struct {
	db *gorm.DB
}

func NewStaticBannerPicRepo() *StaticBannerPicRepo {
	return &StaticBannerPicRepo{
		db: pkgs.GetDB(),
	}
}

const (
	TableStaticBannerPic = "static_banner_pic"
)

func (sr *StaticBannerPicRepo) Get(ctx context.Context, id int64) (*entity.StaticBannerPicBody, error) {
	var app *entity.StaticBannerPicBody
	err := sr.db.Table(TableStaticBannerPic).Where("id = ?", id).First(&app).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return app, nil
}

func (sr *StaticBannerPicRepo) List(ctx context.Context) ([]*entity.StaticBannerPicBody, int64, error) {
	var app []*entity.StaticBannerPicBody
	query := sr.db.Table(TableStaticBannerPic).
		Where("picture != ''")

	err := query.Order("id desc").Find(&app).Error

	if err == gorm.ErrRecordNotFound {
		return nil, 0, nil
	}

	if err != nil {
		return nil, 0, sferror.WithStack(err)
	}

	var total int64
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, sferror.WithStack(err)
	}

	return app, total, nil
}

func (sr *StaticBannerPicRepo) Update(ctx context.Context, id int64, param map[string]interface{}) error {
	err := sr.db.Table(TableStaticBannerPic).Where("id = ?", id).Updates(param).Error
	if err != nil {
		return sferror.WithStack(err)
	}

	return nil
}
