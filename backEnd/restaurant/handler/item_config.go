package handler

import (
	"restaurant/application"
	"restaurant/entity"
	"restaurant/pkgs"

	"github.com/gin-gonic/gin"
)

type ItemConfigHandler struct {
	service application.ItemConfigService
}

func NewItemConfigHandler() *ItemConfigHandler {
	return &ItemConfigHandler{
		service: *application.NewItemConfigService(),
	}
}

func (ih *ItemConfigHandler) Get(ctx *gin.Context) (interface{}, error) {
	return ih.service.Get(ctx)
}

func (ih *ItemConfigHandler) Update(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	var param entity.ItemConfig
	if err := ctx.ShouldBind(&param); err != nil {
		return nil, err
	}

	return nil, ih.service.Update(ctx, userId, param.Prize)
}

func (ih *ItemConfigHandler) Rename(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	var param entity.ItemConfig
	if err := ctx.ShouldBind(&param); err != nil {
		return nil, err
	}

	return nil, ih.service.Rename(ctx, userId, param)
}
