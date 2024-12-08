package handler

import (
	"restaurant/application"
	"restaurant/entity"
	"restaurant/pkgs"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LogoHandler struct {
	service application.LogoService
}

func NewLogoHandler() *LogoHandler {
	return &LogoHandler{
		service: *application.NewLogoService(),
	}
}

func (lh *LogoHandler) Get(ctx *gin.Context) (interface{}, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}

	return lh.service.Get(ctx, id)
}

func (lh *LogoHandler) Update(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))

	if err != nil {
		return nil, err
	}

	var param *entity.Logo
	if err = ctx.ShouldBind(&param); err != nil {
		return nil, err
	}

	return nil, lh.service.Update(ctx, userId, param)
}

func (lh *LogoHandler) List(ctx *gin.Context) (interface{}, error) {
	return lh.service.List(ctx)
}
