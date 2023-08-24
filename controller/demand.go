package controller

import (
	"point-manage/dao"
	"point-manage/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Demand demand

type demand struct{}

//更新需求
func (d *demand) Update(context *gin.Context) {
	params := new(struct {
		Id   uint   `json:"id"`
		Name string `json:"name"`
	})
	if err := context.ShouldBindJSON(params); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}
	err := dao.Demand.Update(&model.Demand{
		ID:         params.Id,
		DemandName: params.Name,
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": nil,
		"code": 200,
	})
}
