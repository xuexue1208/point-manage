package model

//需求表
type Demand struct {
	ID uint `json:"id" gorm:"primary_key"`

	DemandName  string `json:"name" gorm:"column:name"` //需求名称
	VersionCode uint   `json:"versioncode" gorm:"column:versioncode"`
}

func (*Demand) TableName() string {
	return "demand"
}
