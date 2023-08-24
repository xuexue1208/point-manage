package dao

import (
	"point-manage/db"
	"point-manage/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

var Kernel kernel

type kernel struct{}

func (*kernel) List() ([]*model.Kernel, error) {
	var KernelList []*model.Kernel
	db := db.GetPointDB()
	err := db.Model(&model.Kernel{}).Order("id desc").Find(&KernelList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取核心埋点列表失败,%v\n", err))
	}
	return KernelList, nil
}
