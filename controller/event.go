package controller

import (
	"point-manage/dao"
	"point-manage/model"
	"point-manage/service"
	"point-manage/utils"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"net/http"
)

var Event event

type event struct{}

func (*event) Create(ctx *gin.Context) {
	//定义前端传入数据结构
	params := new(struct {
		model.Event
		CategoryList []uint `json:"categoryList"`
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
	err, id := service.Event.Create(&model.Event{
		Event:       params.Event.Event,
		Name:        params.Name,
		Versioncode: params.Versioncode,
		ReportDesc:  params.ReportDesc,
		Remark:      params.Remark,
		CreatedTime: utils.GetNowTimestamp(),
		UpdatedTime: utils.GetNowTimestamp(),
		Imgs:        params.Imgs,
		DemandId:    params.DemandId,
		Operator:    ctx.GetString("username"),
		Status:      0,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}
	//将CategoryList 和返回去的eventid 写入 classevent表
	for _, v := range params.CategoryList {
		err, _ := service.ClassEvent.Create(&model.ClassEvent{
			ClassificationId: v,
			EventId:          id,
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
		"msg":  "新增事件成功,id为" + string(id),
		"data": nil,
		"code": 200,
	})
}

//根据id获取事件详情
func (*event) Info(ctx *gin.Context) {
	params := new(struct {
		EventId uint `json:"eventId"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return

	}
	data, err := service.Event.GetEventById(params.EventId)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取事件详情成功",
		"data": data,
		"code": 200,
	})
}

func (*event) ListByDemand(ctx *gin.Context) {
	params := new(struct {
		DemandId uint `json:"demandId"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return

	}
	data, err := service.Event.ListByDemandId(params.DemandId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取需求事件成功",
		"data": data,
		"code": 200,
	})
}

//根据分类id获取简明的事件列表
func (*event) ConciseListByClassId(ctx *gin.Context) {
	params := new(struct {
		CategoryId uint `json:"categoryId"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return

	}
	data, err := service.Event.ConciseListByCategoryId(params.CategoryId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取需求事件成功",
		"data": data,
		"code": 200,
	})
}

//根据分类id获取事件列表
func (*event) ListByCategoryId(ctx *gin.Context) {
	params := new(struct {
		CategoryId uint `json:"categoryId"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return

	}
	data, err := service.Event.ListByCategoryId(params.CategoryId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取需求事件成功",
		"data": data,
		"code": 200,
	})
}

//存量事件更新版本
func (*event) EventAddVersion(ctx *gin.Context) {
	//定义前端传入数据结构
	params := new(struct {
		Versioncode uint   `json:"versioncode" `
		DemandId    uint   `json:"demandId" `
		EventList   []uint `json:"eventList"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}

	for _, v := range params.EventList {
		info, err := dao.Event.GetEventById(v)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg":  err.Error(),
				"data": nil,
				"code": 403,
			})
			return
		}
		logger.Info(ctx.GetString("username"))
		err = service.Event.UpdatVersion(&model.Event{
			ID:          info.ID,
			Versioncode: params.Versioncode,
			UpdatedTime: utils.GetNowTimestamp(),
			Operator:    ctx.GetString("username"),
			DemandId:    params.DemandId,
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
		"msg":  "更新事件成功",
		"data": nil,
		"code": 200,
	})

}

//更新事件
func (*event) Update(ctx *gin.Context) {
	//定义前端传入数据结构
	params := new(service.ResponseEventDetail)
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}
	if !utils.CheckEn(params.Event.Event) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  "事件名称只能是大小写英文&&数字&&_",
			"data": nil,
			"code": 1002,
		})
		return
	}
	op := ctx.GetString("username")
	//更新事件
	err := service.Event.Update(*params, op)
	if err != nil {
		if err.Error() == "HaveEventNameEn" {
			ctx.JSON(http.StatusOK, gin.H{
				"msg":  "事件名称已存在",
				"data": nil,
				"code": 1001,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新事件成功" + op,
		"data": nil,
		"code": 200,
	})
}

//
type EventParams struct {
	EventId uint              `json:"eventId"`
	Props   []AttributeParams `json:"props"`
}

type AttributeParams struct {
	model.Attribute
	Values []model.Value `json:"values"`
}

//软删除u
func (*event) Delete(ctx *gin.Context) {
	params := new(struct {
		EventIds []uint `json:"eventIds"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}
	for _, v := range params.EventIds {
		err := dao.Event.UpdateEvent(&model.Event{
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

// 根据关键字搜索事件
func (*event) SearchByKeyword(ctx *gin.Context) {
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

	events, err := service.Event.ListBySearchKey(params.Keyword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "搜索成功",
		"data": events,
		"code": 200,
	})

}

// 根据versioncode获取事件列表
func (*event) ListByVersionCode(ctx *gin.Context) {
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

	events, err := service.Event.ListByVersionCode(params.Versioncode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取事件列表成功",
		"data": events,
		"code": 200,
	})

}

//硬删除
func (*event) RealDelete(ctx *gin.Context) {
	params := new(struct {
		EventIds []uint `json:"eventIds"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}
	for _, v := range params.EventIds {
		err := service.Event.RealDelete(v)
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

//Remove
func (*event) Remove(ctx *gin.Context) {
	params := new(struct {
		EventId  uint `json:"eventId"`
		DemandId uint `json:"demandId"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 400,
		})
		return
	}
	err := service.Event.RemoveDemand(params.EventId, params.DemandId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  err.Error(),
			"data": nil,
			"code": 403,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Remove success",
		"data": nil,
		"code": 200,
	})

}
