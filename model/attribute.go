package model

//属性表
type Attribute struct {
	ID uint `json:"id" gorm:"primary_key"`
	//CreatedAt time.Time
	//UpdatedAt time.Time  //更新时间
	//DeletedAt *time.Time `sql:"index"`

	Name        string `json:"name" gorm:"column:name"`               //属性 中文名称
	Key         string `json:"key" gorm:"column:key"`                 //属性 英文名称
	Type        string `json:"type" gorm:"column:type"`               //数据类型
	DemandId    uint   `json:"demandId" gorm:"column:demandId"`       //需求id
	VersionCode uint   `json:"versioncode" gorm:"column:versioncode"` //版本号
	Operator    string `json:"operator" gorm:"column:operator"`       //操作人
	CreatedTime int64  `json:"createdTime" gorm:"column:createdTime"` //创建时间
	UpdatedTime int64  `json:"updatedTime" gorm:"column:updatedTime"` //更新时间
	Remark      string `json:"remark" gorm:"column:remark"`           //事件备注
	Imgs        Array  `json:"imgs" gorm:"column:imgs" `              //配图
	EventId     uint   `json:"event_id"  gorm:"column:event_id"`      //事件id 英文名称
	Status      uint   `json:"status" gorm:"column:status"`           //事件状态 0:正常 1:已删除
}

func (*Attribute) TableName() string {
	return "attribute"
}
