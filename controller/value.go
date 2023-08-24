package controller

import (
	"point-manage/dao"
	"point-manage/model"
	"point-manage/service"
	"point-manage/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Value value

type value struct{}

func (*value) Delete(ctx *gin.Context) {
	params := new(struct {
		ValueIds []uint `json:"valueIds"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}
	for _, v := range params.ValueIds {
		err := dao.Value.DeleteById(&model.Value{
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

//更新
func (*value) Update(ctx *gin.Context) {
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
	err := dao.Value.UpdateStatus(&model.Value{
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

//RealDelete
func (*value) RealDelete(ctx *gin.Context) {
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
	err := service.Value.RealDelete(params.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除成功",
		"data": nil,
		"code": 200,
	})

}

func (*value) ByAttKeyGetValues(ctx *gin.Context) {
	params := new(struct {
		Prop string `json:"prop"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}
	data, err := service.Value.ByAttKeyGetValues(params.Prop)
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
