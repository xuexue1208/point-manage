package dao

import (
	"point-manage/db"
	"point-manage/model"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"gorm.io/gorm"
)

var Demand demand

type demand struct{}

//创建
func (*demand) Create(demand *model.Demand) error {
	db := db.GetPointDB()
	tx := db.Create(&demand)

	if tx.Error != nil {
		return errors.New(fmt.Sprintf("添加需求失败, %v\n", tx.Error))
	}
	return nil

}

//基于DemandName查询demand 表
func (*demand) GetDemandByName(demandName string, versioncode uint) (bool, error) {
	var demand model.Demand
	//数据库查询
	db := db.GetPointDB()
	err := db.Model(&model.Demand{}).Where("name = ? and versioncode = ?", demandName, versioncode).First(&demand).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("查询需求失败,%v\n", err))
		return false, errors.New(fmt.Sprintf("查询需求失败,%v\n", err))
	}
	return false, nil
}

//
func (*demand) VersionDemandList(versioncode uint) ([]*model.Demand, error) {
	var versionDemandList []*model.Demand

	//数据库查询，Limit方法用于限制条数，Offset方法设置起始位置
	db := db.GetPointDB()
	err := db.Model(&model.Demand{}).Where("versioncode = ? ", versioncode).Order("id desc").Find(&versionDemandList).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		logger.Error(fmt.Sprintf("获取版本列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("获取版本列表失败,%v\n", err))
	}

	return versionDemandList, nil
}

//更新
func (*demand) Update(demand *model.Demand) error {
	db := db.GetPointDB()
	tx := db.Model(&model.Demand{}).Where("id = ?", demand.ID).Updates(&demand)
	if tx.Error != nil {
		return errors.New(fmt.Sprintf("更新需求失败, %v\n", tx.Error))
	}
	return nil
}
