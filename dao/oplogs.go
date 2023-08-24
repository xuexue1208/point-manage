package dao

import (
	"point-manage/db"
	"point-manage/model"
)

var Oplogs oplogs

type oplogs struct{}

//创建操作日志
func (*oplogs) Create(oplogs *model.Oplogs) (uint, error) {
	db := db.GetPointDB()
	tx := db.Create(&oplogs)

	if tx.Error != nil {
		return 0, tx.Error
	}
	return oplogs.ID, nil
}
