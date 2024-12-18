package handler

import (
	"restaurant/application"
	"restaurant/entity"
	"restaurant/pkgs"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GuestTotalHandler struct {
	service     application.GuestService
	listService application.GuestShowListService
}

// NewGuestHandler ...
func NewGuestTotalHandler() *GuestTotalHandler {
	return &GuestTotalHandler{
		service:     *application.NewGuestService(),
		listService: *application.NewGuestShowListService(),
	}
}

func (ch *GuestTotalHandler) Create(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	var param entity.GuestResp
	if err := ctx.ShouldBind(&param); err != nil {
		return "", err
	}

	app, err := ch.service.Create(ctx, param, userId)
	if err != nil {
		return nil, err
	}

	data := entity.GuestShowListBody{
		AppId: app.Id,
	}

	_, err = ch.listService.Create(ctx, userId, &data)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (ch *GuestTotalHandler) Delete(ctx *gin.Context) (interface{}, error) {
	accountID, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}

	err = ch.service.Delete(ctx, accountID, id)
	if err != nil {
		return nil, err
	}

	err = ch.listService.Delete(ctx, accountID, id)
	return nil, err
}

func (ch *GuestTotalHandler) List(ctx *gin.Context) (interface{}, error) {
	lang := ctx.DefaultQuery("lang", "zh")
	return ch.service.List(ctx, lang)
}
