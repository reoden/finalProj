package handler

import (
	"restaurant/application"
	"restaurant/entity"
	"restaurant/pkgs"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BannerHandler struct {
	service application.BannerService
}

// NewBannerHandler ...
func NewBannerHandler() *BannerHandler {
	return &BannerHandler{
		service: *application.NewBannerService(),
	}
}

func (bh *BannerHandler) Create(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	var param entity.BannerBody

	if err := ctx.ShouldBind(&param); err != nil {
		return "", err
	}

	return bh.service.Create(ctx, userId, &param)
}

func (bh *BannerHandler) Delete(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	return nil, bh.service.Delete(ctx, userId, id)
}

func (bh *BannerHandler) Update(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	type updateParam struct {
		AppIds []int64 `json:"app_ids"`
	}

	var params updateParam
	if err := ctx.ShouldBind(&params); err != nil {
		return "", err
	}

	return nil, bh.service.Update(ctx, userId, params.AppIds)
}

func (bh *BannerHandler) List(ctx *gin.Context) (interface{}, error) {
	return bh.service.List(ctx)
}
