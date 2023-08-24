package model

//埋点上报表
type Point struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Versioncode uint   `json:"versioncode" gorm:"column:versioncode"`   //版本号
	Event       string `json:"event" gorm:"column:event"`               //事件
	CreatedTime int64  `json:"created_time" gorm:"column:created_time"` //操作时间
}

func (*Point) TableName() string {
	return "point"
}
