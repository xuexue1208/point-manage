package service

import (
	"point-manage/dao"
	"point-manage/model"
	"errors"
	"github.com/wonderivan/logger"
)

var Version version

type version struct{}

//创建
func (*version) Create(version *model.Version) (error, uint) {
	//判断用户名是否存在
	has, err := dao.Version.IfVersionCode(version.VersionCode)
	if err != nil {
		return err, 0
	}
	if has {
		return errors.New("该版本号VersionCode已存在，请重新添加"), 0
	}

	//执行添加用户动作
	if err := dao.Version.Create(version); err != nil {
		return err, 0
	}
	return nil, version.ID
}

//更新
func (*version) Update(version *model.Version) error {
	if err := dao.Version.Update(version); err != nil {
		return err
	}
	//根据versioncode 获取event list
	eventList, err := dao.Event.ListByVersionCode(version.VersionCode)
	if err != nil {
		return err
	}
	ids := make([]uint, 0)
	for _, event := range eventList {
		ids = append(ids, event.ID)
	}
	//更新event
	if err := dao.Event.UpdateEventList(ids); err != nil {
		return err
	}
	return nil
}

//版本详情对应的结构体
type ResponseVersionDetail struct {
	model.Version
	DemandList []*model.Demand
}

//版本详情
func (*version) GetVersionDetail(versioncode uint) (*ResponseVersionDetail, error) {
	//查询版本表
	versioninfo, has, err := dao.Version.GetVersionDetail(versioncode)
	if !has {
		return nil, err
	}
	//查询需求表
	demandList, err := dao.Demand.VersionDemandList(versioncode)
	if err != nil {
		return nil, err
	}

	if len(demandList) == 0 {
		logger.Warn("获取版本下的demandList为空,请确认!")
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &ResponseVersionDetail{
		Version: model.Version{
			VersionName: versioninfo.VersionName,
			VersionCode: versioninfo.VersionCode,
			Operator:    versioninfo.Operator,
			CreatedTime: versioninfo.CreatedTime,
		},
		DemandList: demandList,
	}, nil
}
