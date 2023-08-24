package model

//属性取值表
type Oplogs struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Username    string `json:"username" gorm:"column:username"`         //操作人员
	Method      string `json:"method" gorm:"column:method"`             //操作方法
	Url         string `json:"url" gorm:"column:url"`                   //操作url
	Ip          string `json:"ip" gorm:"column:ip"`                     //操作ip
	Request     string `json:"request" gorm:"column:request"`           //操作数据
	Response    string `json:"response" gorm:"column:response"`         //操作返回数据
	CreatedTime string `json:"created_time" gorm:"column:created_time"` //操作时间
}

func (*Oplogs) TableName() string {
	return "oplogs"
}
