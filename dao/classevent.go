package dao

import (
	"point-manage/db"
	"point-manage/model"
)

var Classevent classevent

type classevent struct{}

//根据eventid 查询所述分类
func (*classevent) ClassIdByEventId(eventid uint) ([]uint, error) {
	var ClasseventList []*model.ClassEvent
	db := db.GetPointDB()
	err := db.Model(&model.ClassEvent{}).Where("event_id = ? ", eventid).Order("id asc").Find(&ClasseventList).Error
	if err != nil {
		return nil, err
	}

	var classIds []uint
	for _, v := range ClasseventList {
		if v.ClassificationId == 0 {
			continue
		}
		classIds = append(classIds, v.ClassificationId)
	}
	return classIds, nil
}

//创建分类
func (*classevent) Create(classevent *model.ClassEvent) error {
	db := db.GetPointDB()
	tx := db.Create(&classevent)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

//根据event_id 删除
func (*classevent) DeleteByEventId(eventid uint) error {
	db := db.GetPointDB()
	tx := db.Where("event_id = ?", eventid).Delete(&model.ClassEvent{})

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
