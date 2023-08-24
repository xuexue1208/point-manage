package dao

import (
	"point-manage/db"
	"point-manage/model"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"gorm.io/gorm"
)

var Value value

type value struct{}

//根据单个属性id 获取value 列表
func (*value) ListByAttributeId(attribute uint) ([]*model.Value, error) {
	var ValueList []*model.Value
	db := db.GetPointDB()
	err := db.Model(&model.Value{}).Where("attributeid = ? ", attribute).Order("id asc").Find(&ValueList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("根据属性获取value列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("根据属性获取value列表失败,%v\n", err))
	}

	return ValueList, nil
}

//ListByAttributeIds
func (*value) ListByAttributeIds(ids []uint) ([]*model.Value, error) {
	var ValueList []*model.Value
	db := db.GetPointDB()
	err := db.Model(&model.Value{}).Where("attributeid in ? ", ids).Order("id asc").Find(&ValueList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("根据属性获取value列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("根据属性获取value列表失败,%v\n", err))
	}

	return ValueList, nil
}

//批量查询取值,根据取值id
func (*value) ListByValueIds(valueIds []uint) ([]*model.Value, error) {
	var ValueList []*model.Value
	db := db.GetPointDB()
	err := db.Model(&model.Value{}).Where("id in ? ", valueIds).Order("id asc").Find(&ValueList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("根据取值获取取值列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("根据取值获取取值列表失败,%v\n", err))
	}

	return ValueList, nil
}

//创建
func (*value) Create(value *model.Value) (uint, error) {
	db := db.GetPointDB()
	tx := db.Create(&value)

	if tx.Error != nil {
		return 0, tx.Error
	}
	return value.ID, nil
}

//标记删除
func (*value) DeleteById(value *model.Value) error {
	db := db.GetPointDB()
	tx := db.Model(&model.Value{}).Where("id = ?", value.ID).Updates(value)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

//update
func (*value) Update(value *model.Value) error {
	db := db.GetPointDB()
	tx := db.Model(&model.Value{}).Where("id = ?", value.ID).Updates(value)

	if tx.Error != nil {
		return tx.Error
	}
	return nil

}

func (*value) UpdateStatus(value *model.Value) error {
	db := db.GetPointDB()
	tx := db.Model(&model.Value{}).Where("id = ?", value.ID).Updates(map[string]interface{}{"status": value.Status})

	if tx.Error != nil {
		return tx.Error
	}
	return nil

}

//根据versioncode 获取value
func (*value) ListByVersionCode(versioncode uint) ([]*model.Value, error) {
	var ValueList []*model.Value
	db := db.GetPointDB()
	err := db.Model(&model.Value{}).Where("versioncode = ? ", versioncode).Order("id asc").Find(&ValueList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("根据versioncode获取value列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("根据versioncode获取value列表失败,%v\n", err))
	}

	return ValueList, nil
}

//根据demandid 获取value
func (*value) ListByDemandId(demandId uint) ([]*model.Value, error) {
	var ValueList []*model.Value
	db := db.GetPointDB()
	err := db.Model(&model.Value{}).Where("demandId = ? ", demandId).Order("id asc").Find(&ValueList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("根据demandId获取value列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("根据demandId获取value列表失败,%v\n", err))
	}

	return ValueList, nil
}

//根据ids 删除value
func (*value) DeleteByIds(ids []uint) error {
	db := db.GetPointDB()
	tx := db.Model(&model.Value{}).Where("id in ?", ids).Delete(&model.Value{})

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

//
func (*value) UpdateDemand(eventid, demandid uint) error {
	//根据eventid 获取属性列表
	attributeList, err := Attribute.ListByEventId(eventid)
	if err != nil {
		return err
	}
	attlist := make([]uint, 0)
	for _, v := range attributeList {
		attlist = append(attlist, v.ID)
	}
	//根据属性ids 查询取值列表
	db := db.GetPointDB()
	var ValueList []*model.Value
	err = db.Model(&model.Value{}).Where("attributeid in ?  and demandId = ?  ", attlist, demandid).Order("id asc").Find(&ValueList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	valueids_list := make([]uint, 0)
	for _, v := range ValueList {
		valueids_list = append(valueids_list, v.ID)
	}
	logger.Info(valueids_list)
	//更新value 的demandid
	tx := db.Model(&model.Value{}).Where("id in ?", valueids_list).Updates(map[string]interface{}{"demandId": 0})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (*value) ByAttKeyGetValues(ids []uint) ([]*model.Value, error) {
	var ValueList []*model.Value
	db := db.GetPointDB()
	err := db.Table("value").Select("distinct value,name,remark,versioncode,imgs").Where("attributeid in ? ", ids).Find(&ValueList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("根据属性获取value列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("根据属性获取value列表失败,%v\n", err))
	}

	return ValueList, nil
}
func (*value) EventIdsBySearchKeyword(keyword string) ([]uint, error) {

	var ValueList []*model.Value
	db := db.GetPointDB()
	err := db.Model(&model.Value{}).Where("name like ? or `value` like ?  or `remark` like  ? ", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").Order("id asc").Find(&ValueList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("根据关键字获取属性列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("根据关键字获取属性列表失败,%v\n", err))
	}
	var aIds []uint
	for _, v := range ValueList {
		aIds = append(aIds, v.AttributeId)
	}

	return aIds, nil
}
