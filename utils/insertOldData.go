package utils

import (
	"point-manage/dao"
	"point-manage/model"
	"point-manage/service"
	"encoding/json"
	"fmt"
	"github.com/wonderivan/logger"
	"os"
)

var JsonFile jsonFile

type jsonFile struct{}

//定义用来接受返回的事件和属性的结构体
type ResponseEventDetail1 struct {
	model.Event
	Properties   []*ResponseAttribute1 `json:"properties" `
	CategoryList []uint                `json:"categoryList" `
}

//定义用来接受返回的属性和属性值的结构体
type ResponseAttribute1 struct {
	model.Attribute
	Values []*model.Value `json:"values" `
}

//读取json文件
func (*jsonFile) ReadJsonFile(path string) error {

	//读取json文件
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	//创建json解码器
	decoder := json.NewDecoder(file)
	params := new([]ResponseEventDetail1)

	err = decoder.Decode(&params)
	if err != nil {
		return err
	}
	//定义一个versioncode 切片
	versionlist := make([]uint, 0)

	for _, v := range *params {
		fmt.Println("event:", v.Event)
		err, event_id := service.Event.Create(&model.Event{
			Event:       v.Event.Event,
			Name:        v.Name,
			Versioncode: v.Versioncode,
			ReportDesc:  v.ReportDesc,
			Remark:      v.Remark,
			CreatedTime: v.CreatedTime,
			UpdatedTime: v.UpdatedTime,
			Imgs:        v.Imgs,
			DemandId:    v.DemandId,
			Operator:    "OP",
			Status:      v.Status,
		})
		if err != nil {
			return err
		}
		logger.Info("事件 --- 插入成功", v, "attribute_Id", event_id)
		versionlist = IfNil(v.Versioncode, versionlist)

		//v.CategoryList 数组去重
		classLocList := RemoveDuplicateElement(v.CategoryList)
		for _, c := range classLocList {
			err, classevent_id := service.ClassEvent.Create(&model.ClassEvent{
				ClassificationId: c,
				EventId:          event_id,
			})
			if err != nil {
				return err
			}
			logger.Info("事件class --- 插入成功", v, "classevent_id", classevent_id)

		}

		for _, v1 := range v.Properties {
			fmt.Println("属性", v1.Attribute)
			attribute_Id, err1 := dao.Attribute.Create(&model.Attribute{
				Name:        v1.Name,
				Key:         v1.Key,
				Type:        v1.Type,
				DemandId:    v1.DemandId,
				VersionCode: v1.VersionCode,
				Operator:    "OP",
				CreatedTime: v1.CreatedTime,
				UpdatedTime: v1.UpdatedTime,
				Remark:      v1.Remark,
				Imgs:        v1.Imgs,
				EventId:     event_id,
				Status:      v1.Status,
			})
			if err1 != nil {
				return err
			}
			logger.Info("属性 --- 插入成功", v1, "attribute_Id", attribute_Id)
			versionlist = IfNil(v1.VersionCode, versionlist)

			for _, v2 := range v1.Values {

				value_id, err := dao.Value.Create(&model.Value{
					Value:       v2.Value,
					Name:        v2.Name,
					Remark:      v2.Remark,
					DemandId:    v2.DemandId,
					Versioncode: v2.Versioncode,
					Imgs:        v2.Imgs,
					Operator:    "OP",
					CreatedTime: v2.CreatedTime,
					UpdatedTime: v2.UpdatedTime,
					AttributeId: attribute_Id,
					Status:      v2.Status,
				})
				if err != nil {
					return err
				}
				logger.Info("value --- 插入成功", v2, "value_id", value_id)
				versionlist = IfNil(v2.Versioncode, versionlist)
			}
		}
	}

	versionList := RemoveDuplicateElement(versionlist)
	logger.Info("versionList", versionList)
	//创建版本
	for _, v := range versionList {
		err, version_id := service.Version.Create(&model.Version{
			VersionCode: v,
			VersionName: string(v),
			Operator:    "OP",
			CreatedTime: GetNowTimestamp(),
			PublishTime: GetNowTimestamp(),
		})
		if err != nil {
			return err
		}
		logger.Info("version --- 插入成功", v, "version_id", version_id)
	}
	return nil
}

func IfNil(v uint, list []uint) []uint {
	if v != 0 {
		list = append(list, v)
		return list
	}
	return list
}
