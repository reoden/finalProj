package handler

import (
	"restaurant/application"
	"restaurant/entity"
	"restaurant/pkgs"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChefHandler struct {
	service application.ChefService
}

// NewChefHandler ...
func NewChefHandler() *ChefHandler {
	return &ChefHandler{
		service: *application.NewChefService(),
	}
}

// Create ...
func (ah *ChefHandler) Create(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	var param entity.ChefResp
	if err := ctx.ShouldBind(&param); err != nil {
		return "", err
	}

	return ah.service.Create(ctx, param, userId)
}

func (ah *ChefHandler) List(ctx *gin.Context) (interface{}, error) {
	lang := ctx.DefaultQuery("lang", "zh")
	return ah.service.List(ctx, lang)
}

// Get ...
func (ah *ChefHandler) Get(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	lang := ctx.DefaultQuery("lang", "zh")

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	return ah.service.Get(ctx, userId, id, lang)
}

// Get ...
func (ah *ChefHandler) Read(ctx *gin.Context) (interface{}, error) {
	lang := ctx.DefaultQuery("lang", "zh")
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	return ah.service.ReadApp(ctx, id, lang)
}

// Delete ...
func (ah *ChefHandler) Delete(ctx *gin.Context) (interface{}, error) {
	accountID, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	return nil, ah.service.Delete(ctx, accountID, id)
}
