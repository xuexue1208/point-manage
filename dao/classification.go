package dao

import (
	"point-manage/db"
	"point-manage/model"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"gorm.io/gorm"
)

var Classification classification

type classification struct{}

//list 获取分类列表
func (*classification) List() ([]*model.Classification, error) {
	var ClassList []*model.Classification
	db := db.GetPointDB()
	err := db.Model(&model.Classification{}).Order("id desc").Find(&ClassList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("获取分类列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("获取分类列表失败,%v\n", err))
	}

	return ClassList, nil
}

//创建分类
func (*classification) Create(classification *model.Classification) error {
	db := db.GetPointDB()
	tx := db.Create(&classification)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

//更新分类
func (*classification) Update(classification *model.Classification) error {
	db := db.GetPointDB()
	tx := db.Model(&model.Classification{}).Where("id = ?", classification.ID).Updates(classification)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

//根据eventid删除分类
func (*classification) DeleteByEventId(eventId uint) error {
	db := db.GetPointDB()
	tx := db.Where("event_id = ?", eventId).Delete(&model.ClassEvent{})

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

//获取指定分类下的事件
func (*classification) ListByClassId(classId uint) ([]*model.ClassEvent, error) {
	var classEventList []*model.ClassEvent
	db := db.GetPointDB()
	err := db.Model(&model.ClassEvent{}).Where("classification_id = ? ", classId).Order("id asc").Find(&classEventList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("获取分类事件失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("获取分类事件失败,%v\n", err))
	}

	return classEventList, nil
}
