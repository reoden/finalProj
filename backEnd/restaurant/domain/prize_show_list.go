package domain

import (
	"context"
	"fmt"
	"restaurant/entity"
	"restaurant/pkgs"
	"restaurant/pkgs/sferror"

	"gorm.io/gorm"
)

const (
	TablePrizeShowList  = "prize_show_list"
	TablePrizeBannerPic = "prize_banner_pic"
)

type PrizeShowListRepo struct {
	db *gorm.DB
}

func NewPrizeShowListRepo() *PrizeShowListRepo {
	return &PrizeShowListRepo{
		db: pkgs.GetDB(),
	}
}

func (br *PrizeShowListRepo) Create(ctx context.Context, bannerBody *entity.BannerBody) (*entity.BannerBody, error) {
	err := br.db.Table(TablePrizeShowList).Create(&bannerBody).Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return bannerBody, nil
}

func (br *PrizeShowListRepo) CreateBanners(ctx context.Context, bannerBody []*entity.BannerBody) ([]*entity.BannerBody, error) {
	err := br.db.Table(TablePrizeShowList).Create(&bannerBody).Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return bannerBody, nil
}

func (br *PrizeShowListRepo) Delete(ctx context.Context, id int64) error {
	err := br.db.Table(TablePrizeShowList).
		Where("app_id = ?", id).
		Updates(map[string]interface{}{
			"status": entity.BannerDeletedStatus,
		}).
		Error
	if err != nil {
		return sferror.WithStack(err)
	}

	return nil
}

func (br *PrizeShowListRepo) DeleteIds(ctx context.Context, ids []int64) error {
	err := br.db.Table(TablePrizeShowList).
		Where("app_id IN ?", ids).
		Updates(map[string]interface{}{
			"status": entity.BannerDeletedStatus,
		}).
		Error
	if err != nil {
		return sferror.WithStack(err)
	}

	return nil
}

func (br *PrizeShowListRepo) DeleteBanners(ctx context.Context, ids []int64) error {
	err := br.db.Table(TablePrizeShowList).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status": entity.BannerDeletedStatus,
		}).
		Error
	if err != nil {
		return sferror.WithStack(err)
	}

	return nil
}

func (br *PrizeShowListRepo) GetByAppId(ctx context.Context, appId int64) (*entity.BannerBody, error) {
	var param entity.BannerBody
	err := br.db.Table(TablePrizeShowList).
		Where("app_id = ?", appId).
		Where("status = ?", entity.BannerShowStatus).
		First(&param).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return &param, nil
}

func (br *PrizeShowListRepo) List(ctx context.Context) ([]entity.BannerBody, int64, error) {
	var param []entity.BannerBody
	query := br.db.Table(TablePrizeShowList).Where("status = ?", entity.BannerShowStatus)
	err := query.Find(&param).Error

	if err != nil {
		return param, 0, sferror.WithStack(err)
	}

	var total int64
	err = query.Count(&total).Error
	if err != nil {
		return param, 0, sferror.WithStack(err)
	}

	return param, total, nil
}

func (br *PrizeShowListRepo) GetAll(ctx context.Context) ([]entity.BannerBody, error) {
	var param []entity.BannerBody
	err := br.db.Table(TablePrizeShowList).
		Where("status = ?", entity.BannerShowStatus).
		Find(&param).
		Error

	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return param, nil
}

func (br *PrizeShowListRepo) SearchByCity(ctx context.Context, key string) ([]entity.Application, int64, error) {
	var apps []entity.Application

	subQuery := br.db.Table(TablePrizeShowList).Select("app_id").Group("app_id")

	query := br.db.Table(TableApp).
		Where("id IN (?)", subQuery).
		Where("address != ''").
		Where("post_name LIKE ?", fmt.Sprintf("%%%s%%", key)).
		Where("status = ?", entity.StatusAccepted)

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return apps, 0, sferror.WithStack(err)
	}

	err = query.Order("id desc").
		Find(&apps).
		Error

	if err != nil {
		return apps, 0, sferror.WithStack(err)
	}

	return apps, total, err
}

func (br *PrizeShowListRepo) SearchByName(ctx context.Context, name string) ([]entity.Application, int64, error) {
	var apps []entity.Application
	subQuery := br.db.Table(TablePrizeShowList).Select("app_id").Group("app_id")

	query := br.db.Table(TableApp).
		Where("id IN (?)", subQuery).
		Where("address != ''").
		Where("name LIKE ?", fmt.Sprintf("%%%s%%", name)).
		Where("status = ?", entity.StatusAccepted)

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return apps, 0, sferror.WithStack(err)
	}

	err = query.Order("id desc").
		Find(&apps).
		Error

	if err != nil {
		return apps, 0, sferror.WithStack(err)
	}

	return apps, total, err
}

func (br *PrizeShowListRepo) SearchByCityAndName(ctx context.Context, name, city string) ([]entity.Application, int64, error) {
	var apps []entity.Application
	subQuery := br.db.Table(TablePrizeShowList).Select("app_id").Group("app_id")

	query := br.db.Table(TableApp).
		Where("id IN (?)", subQuery).
		Where("address != ''").
		Where("name LIKE ?", fmt.Sprintf("%%%s%%", name)).
		Where("post_name LIKE ?", fmt.Sprintf("%%%s%%", city)).
		Where("status = ?", entity.StatusAccepted)

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return apps, 0, sferror.WithStack(err)
	}

	err = query.Order("id desc").
		Find(&apps).
		Error

	if err != nil {
		return apps, 0, sferror.WithStack(err)
	}

	return apps, total, err
}

func (br *PrizeShowListRepo) SearchById(ctx context.Context, id int64) ([]entity.Application, int64, error) {
	var apps []entity.Application
	subQuery := br.db.Table(TablePrizeShowList).Select("app_id").Group("app_id")

	query := br.db.Table(TableApp).
		Where("id IN (?)", subQuery).
		Where("address != ''").
		Where("id LIKE ?", fmt.Sprintf("%%%d%%", id)).
		Where("status = ?", entity.StatusAccepted)

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return apps, 0, sferror.WithStack(err)
	}

	err = query.Order("id desc").
		Find(&apps).
		Error

	if err != nil {
		return apps, 0, sferror.WithStack(err)
	}

	return apps, total, err
}

func (sr *PrizeShowListRepo) UpdateBanner(ctx context.Context, param map[string]interface{}) error {
	err := sr.db.Table(TablePrizeBannerPic).Where("id = 1").Updates(param).Error
	if err != nil {
		return sferror.WithStack(err)
	}

	return nil
}

func (sr *PrizeShowListRepo) GetBanner(ctx context.Context) (*entity.PrizeBanner, error) {
	var banner entity.PrizeBanner
	err := sr.db.Table(TablePrizeBannerPic).Where("id = 1").First(&banner).Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return &banner, nil
}
