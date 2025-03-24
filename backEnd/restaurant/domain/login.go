package domain

import (
	"context"
	"restaurant/entity"
	"restaurant/pkgs"
	"restaurant/pkgs/sferror"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	TableLogin = "account"
)

type LoginRepo struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewLoginRepo() *LoginRepo {
	return &LoginRepo{
		db:  pkgs.GetDB(),
		rdb: pkgs.GetRedis(),
	}
}

type Login struct {
	ID        int64     `json:"id"`
	Phone     string    `json:"phone" form:"phone" binding:"required"`
	Code      string    `json:"code" form:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ar *LoginRepo) UserLoginInfo(ctx context.Context, phone string) (*Login, error) {
	var login Login
	err := ar.db.Table(TableLogin).
		Where("phone = ?", phone).
		First(&login).
		Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return &login, nil
}

func (ar *LoginRepo) AdminLoginInfo(ctx context.Context, phone string) (*Login, error) {
	var login Login
	err := ar.db.Table(TableLogin).
		Where("phone = ?", phone).
		Where("status = ?", entity.AdministratorStatus).
		First(&login).
		Error
	if err != nil {
		return nil, sferror.WithStack(err)
	}

	return &login, nil
}
