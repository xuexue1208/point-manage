package controller

import (
	"point-manage/dao"
	"point-manage/model"
	"point-manage/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Point point

type point struct{}

//create
func (*point) Create(ctx *gin.Context) {
	//定义前端传入数据结构
	params := new(struct {
		model.Point
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}
	//赋值插入数据
	id, err := dao.Point.Create(&model.Point{
		Versioncode: params.Versioncode,
		Event:       params.Event,
		CreatedTime: utils.GetNowTimestamp(),
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
		"msg":  "新增埋点成功",
		"id":   id,
		"code": 200,
	})
}

//select
func (*point) Select(ctx *gin.Context) {
	//定义前端传入数据结构
	params := new(struct {
		Versioncode uint `json:"versioncode"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}

	points, err := dao.Point.Select(params.Versioncode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "查询埋点成功",
		"data": points,
		"code": 200,
	})
}
