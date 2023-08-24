package dao

import (
	"point-manage/db"
	"point-manage/model"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"gorm.io/gorm"
)

var Role role

type role struct{}

func (*role) RoleAuth(username string) (bool, error) {
	data := &model.Role{}
	db := db.GetPointDB()
	tx := db.Where("username = ? and rolename = ?", username, "data").First(&data)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if tx.Error != nil {
		logger.Error(fmt.Sprintf("查询用户权限失败, %v\n", tx.Error))
		return false, errors.New(fmt.Sprintf("查询用户权限失败, %v\n", tx.Error))
	}

	return true, nil
}

//创建角色
func (*role) Create(role *model.Role) (uint, error) {
	db := db.GetPointDB()
	tx := db.Create(&role)

	if tx.Error != nil {
		return 0, tx.Error
	}
	return role.ID, nil
}
