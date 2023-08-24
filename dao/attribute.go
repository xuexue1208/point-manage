package dao

import (
	"point-manage/db"
	"point-manage/model"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"gorm.io/gorm"
)

var Attribute attribute

type attribute struct{}

//
func (*attribute) ListByEventId(eventid uint) ([]*model.Attribute, error) {
	var AttributeList []*model.Attribute
	db := db.GetPointDB()
	err := db.Model(&model.Attribute{}).Where("event_id = ? ", eventid).Order("id asc").Find(&AttributeList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("根据事件获取属性列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("根据事件获取属性列表失败,%v\n", err))
	}

	return AttributeList, nil
}

//批量查询属性,根据属性id
func (*attribute) ListByAttributeIds(attributeIds []uint) ([]*model.Attribute, error) {
	var AttributeList []*model.Attribute
	db := db.GetPointDB()
	err := db.Model(&model.Attribute{}).Where("id in ? ", attributeIds).Order("id asc").Find(&AttributeList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("根据属性获取属性列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("根据属性获取属性列表失败,%v\n", err))
	}

	return AttributeList, nil
}

//创建属性
func (*attribute) Create(attribute *model.Attribute) (uint, error) {
	db := db.GetPointDB()
	tx := db.Create(&attribute)

	if tx.Error != nil {
		return 0, tx.Error
	}
	return attribute.ID, nil
}

//标记为删除
func (*attribute) DeleteById(attribute *model.Attribute) error {
	db := db.GetPointDB()
	tx := db.Model(&model.Attribute{}).Where("id = ?", attribute.ID).Updates(attribute)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

//关键字查询
func (*attribute) EventIdsBySearchKeyword(keyword string) ([]uint, error) {
	var AttributeList []*model.Attribute
	db := db.GetPointDB()
	err := db.Model(&model.Attribute{}).Where("name like ? or `key` like ?  or `remark` like  ? ", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").Order("id asc").Find(&AttributeList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("根据关键字获取属性列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("根据关键字获取属性列表失败,%v\n", err))
	}
	var eventIds []uint
	for _, v := range AttributeList {
		eventIds = append(eventIds, v.EventId)
	}

	return eventIds, nil

}

//更新
func (*attribute) Update(attribute *model.Attribute) error {
	db := db.GetPointDB()
	tx := db.Model(&model.Attribute{}).Where("id = ?", attribute.ID).Updates(attribute)

	if tx.Error != nil {
		return tx.Error
	}
	return nil

}
func (*attribute) UpdateStatus(attribute *model.Attribute) error {
	db := db.GetPointDB()
	tx := db.Model(&model.Attribute{}).Where("id = ?", attribute.ID).Updates(map[string]interface{}{"status": attribute.Status})

	if tx.Error != nil {
		return tx.Error
	}
	return nil

}

//根据versioncode 获取属性列表
func (*attribute) ListByVersionCode(versionCode uint) ([]*model.Attribute, error) {
	var AttributeList []*model.Attribute
	db := db.GetPointDB()
	err := db.Model(&model.Attribute{}).Where("versioncode = ? ", versionCode).Order("id asc").Find(&AttributeList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("根据versionCode获取属性列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("根据versionCode获取属性列表失败,%v\n", err))
	}

	return AttributeList, nil
}

//根据demandid 获取属性列表
func (*attribute) ListByDemandId(demandId uint) ([]*model.Attribute, error) {
	var AttributeList []*model.Attribute
	db := db.GetPointDB()
	err := db.Model(&model.Attribute{}).Where("demandId = ? ", demandId).Order("id asc").Find(&AttributeList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("根据demandId获取属性列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("根据demandId获取属性列表失败,%v\n", err))
	}

	return AttributeList, nil
}

//根据id 批量删除属性
func (*attribute) DeleteByIds(ids []uint) error {
	db := db.GetPointDB()
	tx := db.Where("id in ?", ids).Delete(&model.Attribute{})

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

//判断是否存在
func (*attribute) Tell(key string) (*model.Attribute, error) {
	data := &model.Attribute{}
	db := db.GetPointDB()
	err := db.Model(&model.Attribute{}).Where(" `key` = ?  ", key).First(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("判断属性是否存在失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("判断属性是否存在失败,%v\n", err))
	}

	return data, nil
}

//更新属性的需求
func (*attribute) UpdateDemand(eventid, demandId uint) error {
	var AttributeList []*model.Attribute
	db := db.GetPointDB()
	err := db.Model(&model.Attribute{}).Where("event_id = ?  and demandId = ?  ", eventid, demandId).Order("id asc").Find(&AttributeList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("根据eventid获取属性列表失败,%v\n", err))
		return errors.New(fmt.Sprintf("根据eventid获取属性列表失败,%v\n", err))
	}
	attlist := make([]uint, 0)
	for _, v := range AttributeList {
		attlist = append(attlist, v.ID)
	}
	tx := db.Model(&model.Attribute{}).Where("id in ?", attlist).Updates(map[string]interface{}{"demandId": 0})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

//关键字查询
func (*attribute) BySearchKeywordGetIds(keyword string) ([]uint, error) {
	var AttributeList []*model.Attribute
	db := db.GetPointDB()
	err := db.Model(&model.Attribute{}).Where(" `key` = ? ", keyword).Order("id asc").Find(&AttributeList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("根据关键字获取属性列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("根据关键字获取属性列表失败,%v\n", err))
	}
	var AttIds []uint
	for _, v := range AttributeList {
		AttIds = append(AttIds, v.ID)
	}

	return AttIds, nil

}

//通过key 查询name,type,key,remark 并去重
func (*attribute) Recommend(keyword string) ([]*model.Attribute, error) {
	var AttributeList []*model.Attribute
	db := db.GetPointDB()
	err := db.Table("attribute").Select("distinct name,`type`,`key`,remark").Where(" `key` like ?", "%"+keyword+"%").Find(&AttributeList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintf("根据keyword获取属性列表失败,%v\n", err))
	}
	return AttributeList, nil
}

func (*attribute) GetEventidByaids(ids []uint) ([]uint, error) {
	var AttributeList []*model.Attribute
	db := db.GetPointDB()
	err := db.Model(&model.Attribute{}).Where(" `id` in ? ", ids).Order("id asc").Find(&AttributeList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("根据ids获取属性列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("根据ids获取属性列表失败,%v\n", err))
	}
	var EventIds []uint
	for _, v := range AttributeList {
		EventIds = append(EventIds, v.EventId)
	}

	return EventIds, nil
}
