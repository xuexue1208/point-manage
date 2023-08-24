package service

import (
	"point-manage/dao"
	"point-manage/model"
)

var Demand demand

type demand struct{}

//创建
func (*demand) Create(demand *model.Demand) (error, uint) {
	//执行添加用户动作
	if err := dao.Demand.Create(demand); err != nil {
		return err, 0
	}
	return nil, demand.ID
}

//更新
func (*demand) Update(versioncode uint, demandlist []string) error {
	for _, v := range demandlist {
		bool, err := dao.Demand.GetDemandByName(v, versioncode)

		if err != nil {
			return err
		}
		if bool {
			err := dao.Demand.Create(&model.Demand{DemandName: v, VersionCode: versioncode})
			if err != nil {
				return err
			}
		}
	}
	return nil

}
