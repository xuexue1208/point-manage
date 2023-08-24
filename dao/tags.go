package dao

import (
	"point-manage/db"
	"point-manage/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

var Tags tags

type tags struct{}

//创建
func (*tags) Create(tags *model.Tags) (uint, error) {
	db := db.GetPointDB()
	tx := db.Create(&tags)

	if tx.Error != nil {
		return 0, errors.New(fmt.Sprintf("添加tags失败, %v\n", tx.Error))
	}
	return tags.ID, nil

}

//根据attributeid 获取tags list
func (*tags) ListByAttributeId(AttributeId uint) ([]*model.Tags, error) {
	var tags []*model.Tags
	db := db.GetPointDB()
	tx := db.Where("propertyId = ?", AttributeId).Order("id asc").Find(&tags)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, errors.New(fmt.Sprintf("根据propertyId 获取value list失败, %v\n", tx.Error))
	}
	return tags, nil
}

//根据attributeid 获取tags list
func (*tags) ListByAttributeIds(AttributeIds []uint) ([]*model.Tags, error) {
	var tags []*model.Tags
	db := db.GetPointDB()
	tx := db.Where("propertyId in ?", AttributeIds).Order("id asc").Find(&tags)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, errors.New(fmt.Sprintf("根据propertyId 获取value list失败, %v\n", tx.Error))
	}
	return tags, nil
}

//根据valueid 查询tags
func (*tags) ListByValueId(ValueId uint) ([]*model.Tags, error) {
	var tags []*model.Tags
	db := db.GetPointDB()
	tx := db.Where("valueId = ?", ValueId).Order("id asc").Find(&tags)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, errors.New(fmt.Sprintf("根据valueid 查询tags失败, %v\n", tx.Error))
	}
	return tags, nil
}

//根据valueids 查询tags
func (*tags) ListByValueIds(v_ids []uint) ([]*model.Tags, error) {
	var tags []*model.Tags
	db := db.GetPointDB()
	tx := db.Where("valueId in ?", v_ids).Order("id asc").Find(&tags)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, errors.New(fmt.Sprintf("根据valueid 查询tags失败, %v\n", tx.Error))
	}
	return tags, nil
}

//根据id 删除tag
func (*tags) DeleteById(id uint) error {
	db := db.GetPointDB()
	tx := db.Where("id = ?", id).Delete(&model.Tags{})
	if tx.Error != nil {
		return errors.New(fmt.Sprintf("根据id 删除tag失败, %v\n", tx.Error))
	}
	return nil
}

//根据id 删除tag 批量
func (*tags) DeleteByIds(ids []uint) error {
	db := db.GetPointDB()
	tx := db.Where("id in ?", ids).Delete(&model.Tags{})
	if tx.Error != nil {
		return errors.New(fmt.Sprintf("根据id 删除tag失败, %v\n", tx.Error))
	}
	return nil
}

//根据versioncode 获取tags list
func (*tags) ListByVersionCode(versionCode uint) ([]*model.Tags, error) {
	var tags []*model.Tags
	db := db.GetPointDB()
	tx := db.Where("versioncode = ?", versionCode).Order("id asc").Find(&tags)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, errors.New(fmt.Sprintf("根据versionCode 获取tags list失败, %v\n", tx.Error))
	}
	return tags, nil
}

//根据属性id 删除tags
func (*tags) DeleteByAttributeId(attributeId []uint) error {
	db := db.GetPointDB()
	tx := db.Where("propertyId in ?  and  valueId = 0 and eventId = 0", attributeId).Delete(&model.Tags{})
	if tx.Error != nil {
		return errors.New(fmt.Sprintf("根据属性id 删除tags失败, %v\n", tx.Error))
	}
	return nil
}

//根据valueid 删除tags
func (*tags) DeleteByValueId(valueId []uint) error {
	db := db.GetPointDB()
	tx := db.Where("valueId in ?  and  propertyId = 0  and eventId = 0 ", valueId).Delete(&model.Tags{})
	if tx.Error != nil {
		return errors.New(fmt.Sprintf("根据valueid 删除tags失败, %v\n", tx.Error))
	}
	return nil
}

func (*tags) ListByEventId(eventid uint) ([]*model.Tags, error) {
	var tags []*model.Tags
	db := db.GetPointDB()
	tx := db.Where("eventId = ?  and propertyId = 0 and valueId = 0  ", eventid).Order("id asc").Find(&tags)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if tx.Error != nil {
		return nil, errors.New(fmt.Sprintf("根据versionCode 获取tags list失败, %v\n", tx.Error))
	}
	return tags, nil
}

func (*tags) DeleteEventId(eventid uint) error {
	db := db.GetPointDB()
	tx := db.Where("eventId in ?  and  propertyId = 0  and valueId = 0 ", eventid).Delete(&model.Tags{})
	if tx.Error != nil {
		return errors.New(fmt.Sprintf("根据valueid 删除tags失败, %v\n", tx.Error))
	}
	return nil
}
