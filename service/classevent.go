package service

import (
	"point-manage/dao"
	"point-manage/model"
	"errors"
)

var ClassEvent classEvent

type classEvent struct{}

//创建
func (*classEvent) Create(classEvent *model.ClassEvent) (error, uint) {

	//执行添加class-event动作
	if err := dao.Classevent.Create(classEvent); err != nil {
		return err, 0
	}
	return nil, classEvent.ID
}

//更新
func (*classEvent) Update(classEvent *model.ClassEvent) error {
	//根据eventid 查询是否存在该事件
	eventinfo, err := dao.Event.GetEventById(classEvent.EventId)
	if err != nil {
		return err
	}
	if eventinfo.ID == 0 {
		return errors.New("该事件不存在")
	}

	//执行更新class-event动作
	if err := dao.Classevent.Create(classEvent); err != nil {
		return err
	}
	return nil
}
