package model

//埋点上报表
type Kernel struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Event string `json:"event" gorm:"column:event"` //事件
	Name  string `json:"name" gorm:"column:name"`   //事件 中文名称
}

func (*Kernel) TableName() string {
	return "kernel"
}
