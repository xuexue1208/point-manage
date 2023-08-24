package controller

import (
	"point-manage/dao"
	"point-manage/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Tags tags

type tags struct{}

func (*tags) Create(ctx *gin.Context) {
	//定义前端传入数据结构
	params := new(struct {
		model.Tags
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
	id, err := dao.Tags.Create(&model.Tags{
		Tag:         params.Tag,
		ValueId:     params.ValueId,
		AttributeId: params.AttributeId,
		VersionCode: params.VersionCode,
		EventId:     params.EventId,
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
		"msg":  "新增tag成功,id为" + string(id),
		"data": nil,
		"code": 200,
	})
}

//删除tag
func (*tags) DeleteById(ctx *gin.Context) {
	//定义前端传入数据结构
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
	err := dao.Tags.DeleteById(params.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除tag成功",
		"data": nil,
		"code": 200,
	})
}
