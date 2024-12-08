package handler

import (
	"regexp"
	"restaurant/application"
	"restaurant/entity"
	"restaurant/pkgs"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PrizeShowListHandler struct {
	service application.PrizeShowListService
}

// NewNewsShowListHandler ...
func NewPrizeShowListHandler() *PrizeShowListHandler {
	return &PrizeShowListHandler{
		service: *application.NewPrizeShowListService(),
	}
}

func (bh *PrizeShowListHandler) Create(ctx *gin.Context) (interface{}, error) {
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

func (bh *PrizeShowListHandler) Delete(ctx *gin.Context) (interface{}, error) {
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

func (bh *PrizeShowListHandler) Update(ctx *gin.Context) (interface{}, error) {
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

func (bh *PrizeShowListHandler) List(ctx *gin.Context) (interface{}, error) {
	return bh.service.List(ctx)
}

func (ph *PrizeShowListHandler) Search(ctx *gin.Context) (interface{}, error) {
	key := ctx.DefaultQuery("query", "")
	city := ctx.DefaultQuery("search", "")

	if key == "" && city == "" {
		return ph.List(ctx)
	}
	match, _ := regexp.MatchString("^[0-9]+$", key)

	if !match {
		if city != "" && key == "" {
			ret, err := ph.service.SearchByCity(ctx, city)
			if err != nil {
				return nil, err
			}

			total := ret["total"].(int64)
			if total != 0 {
				return ret, nil
			}
		}
		if city == "" && key != "" {
			return ph.service.SearchByName(ctx, key)
		}

		if city != "" && key != "" {
			return ph.service.SearchByCityAndName(ctx, key, city)
		}

		return nil, nil
	}

	appId, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return nil, err
	}

	return ph.service.SearchById(ctx, appId)

}

func (bh *PrizeShowListHandler) GetPrizeBannerAdver(ctx *gin.Context) (interface{}, error) {
	return bh.service.GetPrizeBannerAdver(ctx), nil
}

func (bh *PrizeShowListHandler) UpdatePrizeBannerAdver(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return "", err
	}

	var param entity.PrizeBanner

	if err = ctx.ShouldBind(&param); err != nil {
		return "", err
	}

	return nil, bh.service.UpdatePrizeBannerAdver(ctx, userId, param)
}
