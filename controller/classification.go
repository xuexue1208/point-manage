package controller

import (
	"point-manage/dao"
	"point-manage/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Classification classification

type classification struct{}

func (*classification) List(ctx *gin.Context) {
	data, err := dao.Classification.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取版本列表成功",
		"data": data,
		"code": 200,
	})
}

//添加分类
func (*classification) Create(ctx *gin.Context) {
	params := new(struct {
		Name string `json:"name" `
		Desc string `json:"desc" `
	})
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  "参数错误",
			"data": nil,
			"code": 403,
		})
		return
	}

	err = dao.Classification.Create(&model.Classification{
		Name: params.Name,
		Desc: params.Desc,
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
		"msg":  "创建分类成功",
		"data": nil,
		"code": 200,
	})
}

//更新分类
func (*classification) Update(ctx *gin.Context) {
	params := new(struct {
		ID   uint   `json:"id" `
		Name string `json:"name"`
		Desc string `json:"desc"`
	})
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  "参数错误",
			"data": nil,
			"code": 403,
		})
		return
	}

	err = dao.Classification.Update(&model.Classification{
		ID:   params.ID,
		Name: params.Name,
		Desc: params.Desc,
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
		"msg":  "更新分类成功",
		"data": nil,
		"code": 200,
	})
}
