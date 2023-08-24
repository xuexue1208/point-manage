package model

//属性取值表
type User struct {
	ID uint `json:"id" gorm:"primary_key"`
	//CreatedAt time.Time
	//UpdatedAt time.Time
	//DeletedAt *time.Time `sql:"index"`

	Name   string `json:"name" gorm:"column:name"`
	Mobile string `json:"mobile" gorm:"column:mobile"`
	Status uint   `json:"status" gorm:"column:status"`
}

func (*User) TableName() string {
	return "xadmin_account"
}
