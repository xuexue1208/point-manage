package service

import (
	"point-manage/dao"
	"point-manage/model"
)

var Value value

type value struct{}

//RealDelete
func (*value) RealDelete(id uint) error {
	//根据value 获取tags
	tags_ids := make([]uint, 0)
	tagsdata_list, err := dao.Tags.ListByValueId(id)
	if err != nil {
		return err
	}
	for _, v := range tagsdata_list {
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
	if err := dao.Value.DeleteByIds([]uint{id}); err != nil {
		return err
	}

	return nil
}

type ByAttKeyGetValuesResponse struct {
	Value       string      `json:"value" gorm:"column:value"`             //取值-EN
	Name        string      `json:"name" gorm:"column:name"`               //取值-CN
	Remark      string      `json:"remark" gorm:"column:remark"`           //取值备注
	Versioncode uint        `json:"versioncode" gorm:"column:versioncode"` //版本号
	Imgs        model.Array `json:"imgs" gorm:"column:imgs"`               //取值配图
}

func (*value) ByAttKeyGetValues(key string) ([]*ByAttKeyGetValuesResponse, error) {
	var ValueList []*ByAttKeyGetValuesResponse
	ids, err := dao.Attribute.BySearchKeywordGetIds(key)
	if err != nil {
		return nil, err
	}
	data, err := dao.Value.ByAttKeyGetValues(ids)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		ValueList = append(ValueList, &ByAttKeyGetValuesResponse{
			Value:       v.Value,
			Name:        v.Name,
			Remark:      v.Remark,
			Versioncode: v.Versioncode,
			Imgs:        v.Imgs,
		})
	}
	return ValueList, nil
}
