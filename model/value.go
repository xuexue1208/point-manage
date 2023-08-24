package model

//属性取值表
type Value struct {
	ID uint `json:"id" gorm:"primary_key"`
	//CreatedAt time.Time
	//UpdatedAt time.Time
	//DeletedAt *time.Time `sql:"index"`

	Value       string `json:"value" gorm:"column:value"`             //取值-EN
	Name        string `json:"name" gorm:"column:name"`               //取值-CN
	Remark      string `json:"remark" gorm:"column:remark"`           //取值备注
	DemandId    uint   `json:"demandId" gorm:"column:demandId"`       //需求id
	Versioncode uint   `json:"versioncode" gorm:"column:versioncode"` //版本号
	Imgs        Array  `json:"imgs" gorm:"column:imgs"`               //取值配图
	Operator    string `json:"operator" gorm:"column:operator"`       //操作人
	CreatedTime int64  `json:"createdTime" gorm:"column:createdTime"` //创建时间
	UpdatedTime int64  `json:"updatedTime" gorm:"column:updatedTime"` //更新时间

	AttributeId uint `json:"attributeid" gorm:"column:attributeid"` //属性id
	Status      uint `json:"status" gorm:"column:status"`           //事件状态 0:正常 1:已删除
}

func (*Value) TableName() string {
	return "value"
}
