package model

//属性取值表
type Role struct {
	ID uint `json:"id" gorm:"primary_key"`

	RoleName string `json:"rolename" gorm:"column:rolename"` //角色名称
	Remark   string `json:"remark" gorm:"column:remark"`     //角色备注
	Mobile   string `json:"mobile" gorm:"column:mobile"`     //角色人员
}

func (*Role) TableName() string {
	return "role"
}
