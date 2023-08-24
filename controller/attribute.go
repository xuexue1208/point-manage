package controller

import (
	"point-manage/dao"
	"point-manage/model"
	"point-manage/service"
	"point-manage/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Attribute attribute

type attribute struct{}

func (*attribute) Delete(ctx *gin.Context) {
	params := new(struct {
		AttributeIds []uint `json:"attributeIds"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}
	for _, v := range params.AttributeIds {
		err := dao.Attribute.DeleteById(&model.Attribute{
			ID:          v,
			UpdatedTime: utils.GetNowTimestamp(),
			Operator:    ctx.GetString("username"),
			Status:      1,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg":  err.Error(),
				"data": nil,
				"code": 403,
			})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除成功",
		"data": nil,
		"code": 200,
	})

}

//判断是否存在EN 属性
func (*attribute) Tell(ctx *gin.Context) {
	params := new(struct {
		Key  string `json:"key"`
		Type string `json:"type"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}
	err, code, data := service.Attribute.Tell(params.Key, params.Type)
	if !utils.CheckEn(params.Key) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  "属性key错误,只能为大小写英文",
			"data": map[string]string{"type": "", "key": ""},
			"code": 2001,
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  err.Error(),
			"data": map[string]string{"type": data.Type, "key": data.Key},
			"code": code,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "校验成功",
		"data": nil,
		"code": 200,
	})
}

func (*attribute) Update(ctx *gin.Context) {
	params := new(struct {
		Id     uint `json:"id"`
		Status uint `json:"status"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}
	err := dao.Attribute.UpdateStatus(&model.Attribute{
		ID:     params.Id,
		Status: params.Status,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": nil,
		"code": 200,
	})

}

//真删除
func (*attribute) RealDelete(ctx *gin.Context) {
	params := new(struct {
		Id uint `json:"id"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}
	err := service.Attribute.RealDelete(params.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": nil,
		"code": 200,
	})

}

//属性推荐
func (*attribute) Recommend(ctx *gin.Context) {
	params := new(struct {
		Keyword string `json:"keyword"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}
	data, err := service.Attribute.Recommend(params.Keyword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"code": 200,
		"data": data,
	})
}
