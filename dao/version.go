package dao

import (
	"point-manage/db"
	"point-manage/model"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"gorm.io/gorm"
)

var Version version

type version struct{}

//创建
func (*version) Create(version *model.Version) error {
	db := db.GetPointDB()
	tx := db.Create(&version)

	if tx.Error != nil {
		return errors.New(fmt.Sprintf("添加版本失败, %v\n", tx.Error))
	}
	return nil

}

//更新
func (*version) Update(version *model.Version) error {
	db := db.GetPointDB()
	tx := db.Model(&model.Version{}).Where("versioncode = ?", version.VersionCode).Updates(version)

	if tx.Error != nil {
		return errors.New(fmt.Sprintf("更新版本失败, %v\n", tx.Error))
	}
	return nil
}

//list
func (*version) List() ([]*model.Version, error) {
	var versionList []*model.Version
	db := db.GetPointDB()
	//数据库查询，Limit方法用于限制条数，Offset方法设置起始位置
	err := db.Model(&model.Version{}).Order("CreatedTime desc").Find(&versionList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("获取版本列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("获取版本列表失败,%v\n", err))
	}

	return versionList, nil
}

//获取版本详情
func (*version) GetVersionDetail(versioncode uint) (*model.Version, bool, error) {
	data := &model.Version{}
	db := db.GetPointDB()
	err := db.Where("versioncode = ?", versioncode).First(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, errors.New("该版本不存在")
	}
	if err != nil {
		logger.Error(fmt.Sprintf("查询版本详情失败, %v\n", err))
		return nil, false, errors.New(fmt.Sprintf("查询版本详情失败, %v\n", err))
	}

	return data, true, nil
}

//判断版本是否纯在
func (*version) IfVersionCode(name uint) (bool, error) {
	data := &model.Version{}
	db := db.GetPointDB()
	err := db.Where("versioncode = ?", name).First(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("查询VersionCode失败, %v\n", err))
		return false, errors.New(fmt.Sprintf("查询VersionCode失败, %v\n", err))
	}

	return true, nil
}
