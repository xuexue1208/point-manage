package service

import (
	"point-manage/dao"
	"point-manage/model"
	"fmt"
)

var Attribute attribute

type attribute struct{}

//判断EN 是否存在
func (*attribute) Tell(key string, t string) (error, uint, *model.Attribute) {
	data, err := dao.Attribute.Tell(key)
	if err != nil {
		return err, 0, nil
	}
	if data == nil {
		return nil, 0, nil
	}
	if data.Key == key && data.Type != t {
		return fmt.Errorf("KEY已经存在,本次添加的与存在的TYPE不一致"), 2003, data
	}
	if data.Key != key && data.Type == t {
		return fmt.Errorf("KEY已经存在,本次添加的与存在的KEY大小写不一致"), 2002, data
	}
	if data.Key != key && data.Type != t {
		return fmt.Errorf("KEY已经存在,本次添加的与存在的KEY/TYPE都不一致"), 2004, data
	}
	if data.Key == key && data.Type == t {
		return nil, 0, data
	}
	return nil, 0, data
}

//真删除
func (*attribute) RealDelete(id uint) error {
	tags_ids := make([]uint, 0)
	//根据属性id 查询tags
	att_tags, err := dao.Tags.ListByAttributeId(id)
	if err != nil {
		return err
	}
	for _, v := range att_tags {
		tags_ids = append(tags_ids, v.ID)
	}
	//根据id查询取值
	value_ids := make([]uint, 0)
	valuedata_list, err := dao.Value.ListByAttributeId(id)
	for _, v := range valuedata_list {
		value_ids = append(value_ids, v.ID)
	}
	//根据value_ids 查询对应的tags
	value_tags, err := dao.Tags.ListByValueIds(value_ids)
	for _, v := range value_tags {
		tags_ids = append(tags_ids, v.ID)
	}

	//删除tags
	if len(tags_ids) > 0 {
		err = dao.Tags.DeleteByIds(tags_ids)
		if err != nil {
			return err
		}
	}
	//删除value
	if len(value_ids) > 0 {
		err = dao.Value.DeleteByIds(value_ids)
		if err != nil {
			return err
		}
	}
	//删除属性
	err = dao.Attribute.DeleteByIds([]uint{id})
	if err != nil {
		return err
	}

	return nil
}

type KeywordResponse struct {
	Name   string `json:"name" gorm:"column:name"`     //属性 中文名称
	Key    string `json:"key" gorm:"column:key"`       //属性 英文名称
	Type   string `json:"type" gorm:"column:type"`     //数据类型
	Remark string `json:"remark" gorm:"column:remark"` //事件备注
}

func (*attribute) Recommend(keyword string) ([]*KeywordResponse, error) {
	var RecommendList []*KeywordResponse
	data, err := dao.Attribute.Recommend(keyword)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		RecommendList = append(RecommendList, &KeywordResponse{
			Name:   v.Name,
			Key:    v.Key,
			Type:   v.Type,
			Remark: v.Remark,
		})
	}
	return RecommendList, nil
}
