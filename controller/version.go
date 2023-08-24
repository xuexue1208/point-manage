package controller

import (
	"point-manage/dao"
	"point-manage/model"
	"point-manage/service"
	"point-manage/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
)

var Version version

type version struct{}

func (*version) Create(ctx *gin.Context) {
	//定义前端传入数据结构
	params := new(struct {
		VersionName string   `json:"versionName" ` //版本名称  225
		VersionCode uint     `json:"versioncode" ` //版本号   20822500
		DemandList  []string `json:"demandList" `  //需求列表

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
	err, id := service.Version.Create(&model.Version{
		VersionName: params.VersionName,
		VersionCode: params.VersionCode,
		Operator:    ctx.GetString("username"),
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
	//将version code 和前端传递的name  写入 demand表
	for _, v := range params.DemandList {
		err, _ := service.Demand.Create(&model.Demand{
			DemandName:  v,
			VersionCode: params.VersionCode,
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
		"msg":  "新增版本成功,id为" + string(id),
		"data": nil,
		"code": 200,
	})
}

func (*version) List(ctx *gin.Context) {

	data, err := dao.Version.List()
	//降序排序
	sort.Slice(data, func(i, j int) bool {
		if data[i].VersionCode > data[j].VersionCode {
			return true
		} else {
			return false
		}
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
		"msg":  "获取版本列表成功",
		"data": data,
		"code": 200,
	})
}

func (*version) VersionInfo(ctx *gin.Context) {
	params := new(struct {
		VersionCode uint `json:"versioncode" `
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}
	data, err := service.Version.GetVersionDetail(params.VersionCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": data,
			"code": 403,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取版本信息成功",
		"code": 200,
		"data": data,
	})
}

func (*version) Update(ctx *gin.Context) {
	params := new(struct {
		VersionCode uint     `json:"versioncode" `
		VersionName string   `json:"versionName" `
		DemandList  []string `json:"demandList" `
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}
	//更新version表
	err := service.Version.Update(&model.Version{
		VersionName: params.VersionName,
		VersionCode: params.VersionCode,
		Operator:    ctx.GetString("username"),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}
	//更新demand表中的相关数据数据,无记录插入,有记录无操作
	err = service.Demand.Update(params.VersionCode, params.DemandList)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "修改版本信息成功",
		"code": 200,
		"data": nil,
	})
}

//更新版本状态
func (*version) UpdateToRelease(ctx *gin.Context) {
	params := new(struct {
		VersionCode uint `json:"versioncode" `
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}
	err := service.Version.Update(&model.Version{
		VersionCode: params.VersionCode,
		Operator:    ctx.GetString("username"),
		PublishTime: utils.GetNowTimestamp(),
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
		"msg":  "修改版本状态成功",
		"code": 200,
		"data": nil,
	})
}
