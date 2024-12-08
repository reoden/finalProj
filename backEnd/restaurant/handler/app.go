package handler

import (
	"errors"
	"io"
	"path"
	"regexp"
	"restaurant/application"
	"restaurant/entity"
	"restaurant/pkgs"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AppHandler struct {
	service application.AppService
}

// NewAppHandler ...
func NewAppHandler() *AppHandler {
	return &AppHandler{
		service: *application.NewAppService(),
	}
}

func UploadPicture(ctx *gin.Context) (interface{}, error) {
	type pairPics struct {
		FileName string `json:"file_name"`
		FileUrl  string `json:"file_url"`
	}

	var pics pairPics

	if pic, err := ctx.FormFile("file"); err == nil {
		//获取文件的后缀名
		extString := path.Ext(pic.Filename)
		//允许上传文件的格式
		allowExtMap := map[string]struct{}{
			".jpg":  {},
			".png":  {},
			".jpeg": {},
		}
		if _, ok := allowExtMap[extString]; !ok {
			return nil, errors.New("上传图片格式错误")
		}

		fileHandler, err := pic.Open()
		if err != nil {
			return nil, errors.New("图像文件打开失败")
		}
		defer fileHandler.Close()

		data, _ := io.ReadAll(fileHandler)
		fileType := strings.Replace(extString, ".", "", -1)
		filename, err := pkgs.UploadAndFilename(ctx, data, fileType, false)
		if err != nil {
			return nil, pkgs.ServerError
		}

		signedUrl, _ := pkgs.SignedUrl(filename, false)
		pics = pairPics{
			FileName: filename,
			FileUrl:  signedUrl,
		}

		return pics, nil
	}

	return nil, errors.New("上传图片失败")
}

// Create ...
func (ah *AppHandler) Create(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	var param entity.AppResp
	if err := ctx.ShouldBind(&param); err != nil {
		return "", err
	}

	return ah.service.Create(ctx, param, userId, entity.StatusPending)
}

// Create ...
func (ah *AppHandler) Save(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	var param entity.AppResp
	if err := ctx.ShouldBind(&param); err != nil {
		return "", err
	}

	return ah.service.Create(ctx, param, userId, entity.StatusSaved)
}

func (ah *AppHandler) List(ctx *gin.Context) (interface{}, error) {
	spage := ctx.DefaultQuery("page", "1")
	ssize := ctx.DefaultQuery("size", "20")
	lang := ctx.DefaultQuery("lang", "zh")

	page, _ := strconv.Atoi(spage)
	size, _ := strconv.Atoi(ssize)
	if page <= 0 || size <= 0 {
		return nil, pkgs.PageError
	}
	return ah.service.List(ctx, (page-1)*size, size, lang)
}

// Get ...
func (ah *AppHandler) Get(ctx *gin.Context) (interface{}, error) {
	accountId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	return ah.service.GetSelfApp(ctx, accountId, id)
}

// Get ...
func (ah *AppHandler) Self(ctx *gin.Context) (interface{}, error) {
	accountId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	return ah.service.GetAppByUserId(ctx, accountId)
}

// Get ...
func (ah *AppHandler) Read(ctx *gin.Context) (interface{}, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}
	return ah.service.ReadApp(ctx, id)
}

func (ah *AppHandler) SearchByCity(ctx *gin.Context) (interface{}, error) {
	key := ctx.DefaultQuery("query", "")
	spage := ctx.DefaultQuery("page", "1")
	ssize := ctx.DefaultQuery("page_size", "20")

	page, _ := strconv.Atoi(spage)
	size, _ := strconv.Atoi(ssize)
	if page <= 0 || size <= 0 {
		return nil, pkgs.PageError
	}

	return ah.service.SearchByCity(ctx, (page-1)*size, size, key)
}

// Delete ...
func (ah *AppHandler) Delete(ctx *gin.Context) (interface{}, error) {
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

func (ah *AppHandler) UpdateStatus(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}

	var param map[string]string

	if err = ctx.ShouldBind(&param); err != nil {
		return nil, err
	}

	status := param["status"]
	return nil, ah.service.UpdateStatus(ctx, userId, id, status)
}

func (ah *AppHandler) Search(ctx *gin.Context) (interface{}, error) {
	key := ctx.DefaultQuery("query", "")
	city := ctx.DefaultQuery("search", "")
	spage := ctx.DefaultQuery("page", "1")
	ssize := ctx.DefaultQuery("page_size", "20")

	page, _ := strconv.Atoi(spage)
	size, _ := strconv.Atoi(ssize)
	if page <= 0 || size <= 0 {
		return nil, pkgs.PageError
	}

	if key == "" && city == "" {
		return ah.SearchList(ctx)
	}
	match, _ := regexp.MatchString("^[0-9]+$", key)

	if !match {
		if city != "" && key == "" {
			ret, err := ah.service.SearchByCity(ctx, (page-1)*size, size, city)
			if err != nil {
				return nil, err
			}

			total := ret["total"].(int64)
			if total != 0 {
				return ret, nil
			}
		}
		if city == "" && key != "" {
			return ah.service.SearchByName(ctx, (page-1)*size, size, key)
		}

		if city != "" && key != "" {
			return ah.service.SearchByCityAndName(ctx, (page-1)*size, size, key, city)
		}

		return nil, nil
	}

	appId, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return nil, err
	}

	return ah.service.SearchById(ctx, (page-1)*size, size, appId)
}

func (ah *AppHandler) SearchList(ctx *gin.Context) (interface{}, error) {
	spage := ctx.DefaultQuery("page", "1")
	ssize := ctx.DefaultQuery("page_size", "20")

	page, _ := strconv.Atoi(spage)
	size, _ := strconv.Atoi(ssize)
	if page <= 0 || size <= 0 {
		return nil, pkgs.PageError
	}
	return ah.service.SearchList(ctx, (page-1)*size, size)
}

func (ah *AppHandler) DownloadDoc(ctx *gin.Context) (interface{}, error) {
	userId, err := pkgs.GetAccountIdFromHeader(ctx.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return nil, err
	}

	return ah.service.DownloadDoc(ctx, userId, id)
}
