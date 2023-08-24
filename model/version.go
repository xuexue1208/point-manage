package model

//版本表
//version: {
//	Demand: {
//		Event: {
//			Attribute: {}
//		}
//	}
//}

type Version struct {
	ID uint `json:"id" gorm:"primary_key"`

	VersionName string `json:"versionName" gorm:"column:versionName"` //版本名称  225
	VersionCode uint   `json:"versioncode" gorm:"column:versioncode"` //版本号   20822500

	Operator    string `json:"operator" gorm:"column:operator"`       //操作人
	PublishTime int64  `json:"publishTime" gorm:"column:publishTime"` //入库（发版）时间
	CreatedTime int64  `json:"createdTime" gorm:"column:createdTime"` //创建时间

}

func (*Version) TableName() string {
	return "version"
}
