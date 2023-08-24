package model

//分类表
type Classification struct {
	ID uint `json:"id" gorm:"primary_key"`
	//CreatedAt time.Time
	//UpdatedAt time.Time
	//DeletedAt *time.Time `sql:"index"`

	Name string `json:"name"  gorm:"column:name"` //分类名称
	Desc string `json:"desc"  gorm:"column:desc"` //分类描述

}

func (*Classification) TableName() string {
	return "classification"
}
