package dao

import (
	"point-manage/db"
	"point-manage/model"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"gorm.io/gorm"
)

var Event event

type event struct{}

func (*event) Create(event *model.Event) error {
	db := db.GetPointDB()
	tx := db.Create(&event)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

//根据demandid 获取event list
func (*event) ListByDemandId(demandId uint) ([]*model.Event, error) {
	var demandEventList []*model.Event

	//数据库查询，Limit方法用于限制条数，Offset方法设置起始位置
	db := db.GetPointDB()
	err := db.Model(&model.Event{}).Where("demandId = ? or tempDemandId = ?", demandId, demandId).Order("id asc").Find(&demandEventList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("获取需求事件失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("获取需求事件失败,%v\n", err))
	}

	return demandEventList, nil
}

//根据eventid 获取event
func (*event) GetEventById(eventId uint) (*model.Event, error) {
	var event model.Event
	db := db.GetPointDB()
	err := db.Model(&model.Event{}).Where("id = ? ", eventId).Order("id asc").Find(&event).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("获取事件失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("获取事件失败,%v\n", err))
	}

	return &event, nil
}

//更新event
func (*event) UpdateEvent(e *model.Event) error {
	logger.Info("更新event", e)
	db := db.GetPointDB()
	//tx := db.GORM.Model(&model.Event{}).Where("id = ?", e.ID).Select("tempDemandId", "status", "demandId", "imgs", "updatedTime", "createdTime", "remark", "reportDesc", "operator", "versioncode", "name", "event").Updates(&e)
	tx := db.Model(&model.Event{}).Where("id = ?", e.ID).Updates(map[string]interface{}{
		"tempDemandId": e.TempDemandId,
		"status":       e.Status,
		"demandId":     e.DemandId,
		"imgs":         e.Imgs,
		"updatedTime":  e.UpdatedTime,
		"createdTime":  e.CreatedTime,
		"remark":       e.Remark,
		"reportDesc":   e.ReportDesc,
		"operator":     e.Operator,
		"versioncode":  e.Versioncode,
		"name":         e.Name,
		"event":        e.Event,
	})

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

//根据ids 批量更新event tempDemandId的值
func (*event) UpdateEventList(ids []uint) error {
	db := db.GetPointDB()
	tx := db.Model(&model.Event{}).Where("id in ?", ids).Updates(map[string]interface{}{"tempDemandId": 0})

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

//根据id 更新event DemandId
func (*event) UpdateEventDemand(id uint) error {
	db := db.GetPointDB()
	tx := db.Model(&model.Event{}).Where("id =  ?", id).Updates(map[string]interface{}{"demandId": 0, "tempDemandId": 0})

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

//根据id 批量查询
func (*event) ListByEventIds(eventIds []uint) ([]*model.Event, error) {
	var eventList []*model.Event
	db := db.GetPointDB()
	err := db.Model(&model.Event{}).Where("id in ? ", eventIds).Order("id asc").Find(&eventList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("根据事件id获取事件列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("根据事件id获取事件列表失败,%v\n", err))
	}

	return eventList, nil
}

//关键字查询
func (*event) EventIdsBySearchKeyword(keyword string) ([]uint, error) {
	var eventList []*model.Event
	db := db.GetPointDB()
	err := db.Model(&model.Event{}).Where("name like ? or `event` like ?  or  `remark` like  ? ", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").Order("id asc").Find(&eventList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("根据关键字查询事件失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("根据关键字查询事件失败,%v\n", err))
	}

	var eventIds []uint
	for _, event := range eventList {
		eventIds = append(eventIds, event.ID)
	}
	return eventIds, nil
}

//根据versioncode 获取event list
func (*event) ListByVersionCode(versionCode uint) ([]*model.Event, error) {
	var eventList []*model.Event
	db := db.GetPointDB()
	err := db.Model(&model.Event{}).Where("versioncode = ? ", versionCode).Order("id asc").Find(&eventList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("获取事件列表失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("获取事件列表失败,%v\n", err))
	}

	return eventList, nil
}

//删除
func (*event) Delete(event *model.Event) error {
	db := db.GetPointDB()
	tx := db.Delete(&event)

	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

//根据event en 查询event
func (*event) GetEventByEnName(eventEn string) (*model.Event, error) {
	var event model.Event
	db := db.GetPointDB()
	err := db.Model(&model.Event{}).Where("event = ? ", eventEn).Order("id asc").Find(&event).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("获取事件失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("获取事件失败,%v\n", err))
	}

	return &event, nil
}

func (*event) ListKernel() ([]*model.Event, error) {
	var eventList []*model.Event
	db := db.GetPointDB()
	err := db.Model(&model.Event{}).Where("kernel = 1  ").Order("id asc").Find(&eventList).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		logger.Error(fmt.Sprintf("获取事件失败,%v\n", err))
		return nil, errors.New(fmt.Sprintf("获取事件失败,%v\n", err))
	}

	return eventList, nil
}

func (*event) TagEventKernel(id, value uint) error {
	db := db.GetPointDB()
	tx := db.Model(&model.Event{}).Where("id = ?", id).Update("kernel", value)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
