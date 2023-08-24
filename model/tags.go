package model

//属性取值表
type Tags struct {
	ID uint `json:"id" gorm:"primary_key"`

	Tag         uint `json:"tag" gorm:"column:tag"`                 //标签 1已有需要检查 2 核心埋点 3标黄
	ValueId     uint `json:"valueId" gorm:"column:valueId"`         //取值id
	AttributeId uint `json:"propertyId" gorm:"column:propertyId"`   //属性id
	EventId     uint `json:"eventId" gorm:"column:eventId"`         //事件id
	VersionCode uint `json:"versioncode" gorm:"column:versioncode"` //版本号
}

func (*Tags) TableName() string {
	return "tags"
}
