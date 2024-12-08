package handler

import (
	"restaurant/application"
	"restaurant/pkgs"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StaticBannerPicHandler struct {
	service application.StaticBannerPicService
}

func NewStaticBannerPicHandler() *StaticBannerPicHandler {
	return &StaticBannerPicHandler{
		service: *application.NewStaticBannerPicService(),
	}
}

func (sh *StaticBannerPicHandler) Get(ctx *gin.Context) (interface{}, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}

	return sh.service.Get(ctx, id)
}

func (sh *StaticBannerPicHandler) List(ctx *gin.Context) (interface{}, error) {
	return sh.service.List(ctx)
}

func (sh *StaticBannerPicHandler) Update(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}
	type staticBannerPicParam struct {
		Picture string `json:"picture"`
		RedirectUrl string `json:"redirect_url"`
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}

	var param staticBannerPicParam
	if err = ctx.ShouldBind(&param); err != nil {
		return nil, err
	}

	return nil, sh.service.Update(ctx, userId, id, param.Picture, param.RedirectUrl)
}
