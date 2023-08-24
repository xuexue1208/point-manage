package dao

import (
	"point-manage/db"
	"point-manage/model"
	"errors"
	"fmt"
)

var Point point

type point struct{}

func (*point) Create(point *model.Point) (uint, error) {
	db := db.GetPointDB()
	tx := db.Create(&point)

	if tx.Error != nil {
		return 0, errors.New(fmt.Sprintf("添加埋点失败, %v\n", tx.Error))
	}
	return point.ID, nil
}

//select
func (*point) Select(versioncode uint) ([]*model.Point, error) {
	var point []*model.Point
	//point := make([]*model.Point, 0)
	db := db.GetPointDB()
	tx := db.Where("versioncode = ?", versioncode).Order("id desc").Find(&point)
	if tx.Error != nil {
		return nil, errors.New(fmt.Sprintf("根据versioncode 获取埋点失败, %v\n", tx.Error))
	}
	return point, nil
}
