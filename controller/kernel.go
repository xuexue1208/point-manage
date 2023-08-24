package controller

import (
	"point-manage/dao"
	"point-manage/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Kernel kernel

type kernel struct{}

//list
func (*kernel) ListKernel(ctx *gin.Context) {
	data, err := service.Event.ListKernel()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取核心埋点列表成功",
		"data": data,
		"code": 200,
	})
}

func (*kernel) TagEventKernalCreate(ctx *gin.Context) {
	params := new(struct {
		EventId uint `json:"eventId"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}
	err := dao.Event.TagEventKernel(params.EventId, 1)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "标记核心埋点成功",
		"code": 200,
	})
}

func (*kernel) TagEventKernalRemove(ctx *gin.Context) {
	params := new(struct {
		EventId uint `json:"eventId"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}
	err := dao.Event.TagEventKernel(params.EventId, 0)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "remove核心埋点成功",
		"code": 200,
	})
}
