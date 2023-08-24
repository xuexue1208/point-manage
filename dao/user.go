package dao

import (
	"point-manage/db"
	"point-manage/model"
	"errors"
	"gorm.io/gorm"
)

var User user

type user struct{}

//根据mobile 查询用户first
func (*user) GetByMobile(mobile string) (*model.User, *model.Role, error, error) {
	var user model.User
	db_x := db.GetXadminDB()
	err := db_x.Model(&model.User{}).Where("mobile = ? and  status = 1", mobile).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err, nil
	}
	if err != nil {
		return nil, nil, err, nil
	}
	var role model.Role
	db_p := db.GetPointDB()
	err = db_p.Model(&model.Role{}).Where("mobile = ? ", mobile).First(&role).Error
	return &user, &role, nil, nil
}
