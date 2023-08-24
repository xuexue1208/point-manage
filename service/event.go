package service

import (
	"point-manage/dao"
	"point-manage/model"
	"fmt"
	"github.com/wonderivan/logger"
	"sort"
	"sync"
	"time"
)

var Event event

type event struct{}

//创建
func (*event) Create(event *model.Event) (error, uint) {
	//根据event en查询是否存在
	eventinfo, err := dao.Event.GetEventByEnName(event.Event)
	if err != nil {
		return err, 0
	}
	if eventinfo.ID != 0 {
		return fmt.Errorf("该事件En已存在,更换名称"), 0
	}
	//执行添加事件动作
	if err := dao.Event.Create(event); err != nil {
		return err, 0
	}
	return nil, event.ID
}

//定义用来接受返回的事件和属性的结构体
type ResponseEventDetail struct {
	model.Event
	Properties   []*ResponseAttribute `json:"properties" `
	CategoryList []uint               `json:"categoryList" `
	Tags         []*model.Tags        `json:"tags" `
}

//定义用来接受返回的属性和属性值的结构体
type ResponseAttribute struct {
	model.Attribute
	Values []*ResponseValue `json:"values" `
	Tags   []*model.Tags    `json:"tags" `
}

type ResponseValue struct {
	model.Value
	Tags []*model.Tags `json:"tags" `
}

type ResponseEventInfo struct {
	ID    uint   `json:"eventId" gorm:"primary_key"`
	Event string `json:"event" gorm:"column:event"` //事件 英文名称
	Name  string `json:"name" gorm:"column:name"`   //事件 中文名称
}

//根据categoryId 获取简明的event list
func (*event) ConciseListByCategoryId(categoryId uint) ([]*ResponseEventInfo, error) {
	//根据demandid 获取event list
	eventlist, err := dao.Classification.ListByClassId(categoryId)
	if err != nil {
		return nil, err
	}
	eventids := make([]uint, 0)
	for _, k := range eventlist {
		eventids = append(eventids, k.EventId)
	}

	events_data, err := dao.Event.ListByEventIds(eventids)
	if err != nil {
		return nil, err
	}
	//定义一个事件list,用来存放事件和属性
	responseEventList := make([]*ResponseEventInfo, 0)
	var wg sync.WaitGroup
	var lock sync.Mutex
	for _, eventinfo := range events_data {
		wg.Add(1)
		go func(eventinfo *model.Event) {
			defer wg.Done()
			if err != nil {
				return //nil, err
			}
			info := &ResponseEventInfo{
				ID:    eventinfo.ID,
				Event: eventinfo.Event,
				Name:  eventinfo.Name,
			}
			lock.Lock()
			responseEventList = append(responseEventList, info)
			defer lock.Unlock()
		}(eventinfo)
		wg.Wait()

	}

	return responseEventList, nil
}

//根据categoryId 获取event list
func (*event) ListByCategoryId(categoryId uint) ([]*ResponseEventDetail, error) {
	//根据demandid 获取event list
	classeventlist, err := dao.Classification.ListByClassId(categoryId)
	if err != nil {
		return nil, err
	}
	//定义一个切片，用来存放eventid
	eventIdList := make([]uint, 0)
	for _, k := range classeventlist {
		eventIdList = append(eventIdList, k.EventId)
	}
	//批量查询event
	events_data, err := dao.Event.ListByEventIds(eventIdList)
	var wg sync.WaitGroup
	var lock sync.Mutex
	//定义一个事件list,用来存放事件和属性
	responseEventList := make([]*ResponseEventDetail, 0)

	for _, event := range events_data {
		wg.Add(1)
		go func(event *model.Event) {
			defer wg.Done()
			info, err := Event.GetEventById_Go(event)
			if err != nil {
				return //nil, err
			}
			lock.Lock()
			defer lock.Unlock()
			responseEventList = append(responseEventList, info)
		}(event)
	}
	wg.Wait()
	//升序排序
	sort.Slice(responseEventList, func(i, j int) bool {
		if responseEventList[i].Event.ID < responseEventList[j].Event.ID {
			return true
		} else {
			return false
		}
	})
	return responseEventList, nil
}

//存量事件更新版本
func (*event) UpdatVersion(event *model.Event) error {
	//执行更新事件动作
	if err := dao.Event.UpdateEvent(event); err != nil {
		return err
	}
	return nil
}

//更新event
func (*event) Update(info ResponseEventDetail, op string) error {
	//如果事件id 为0 , 则为新增事件 ,包括新增 属性 取值
	NowTime := time.Now().Unix() * 1000

	//info.id 等于0 新增,不等于0 更新
	if info.ID == 0 {
		err, event_id := Event.Create(&model.Event{
			Event:        info.Event.Event,
			Name:         info.Name,
			Versioncode:  info.Versioncode,
			Operator:     op,
			ReportDesc:   info.ReportDesc,
			Remark:       info.Remark,
			CreatedTime:  NowTime,
			UpdatedTime:  NowTime,
			Imgs:         info.Imgs,
			DemandId:     info.DemandId,
			Status:       0,
			TempDemandId: 0,
		})
		if err != nil {
			return fmt.Errorf("HaveEventNameEn")
		}
		//将CategoryList 和返回去的eventid 写入 classevent表
		for _, v := range info.CategoryList {
			err, _ := ClassEvent.Create(&model.ClassEvent{
				ClassificationId: v,
				EventId:          event_id,
			})
			if err != nil {
				return err
			}
		}
		//将AttributeList 和返回去的eventid 写入 attribute表
		for _, v := range info.Properties {
			att_id, err := dao.Attribute.Create(&model.Attribute{
				Name:        v.Name,
				Key:         v.Key,
				Type:        v.Type,
				DemandId:    v.DemandId,
				VersionCode: v.VersionCode,
				Operator:    op,
				CreatedTime: NowTime,
				UpdatedTime: NowTime,
				Remark:      v.Remark,
				Imgs:        v.Imgs,
				EventId:     event_id,
				Status:      0,
			})
			if err != nil {
				return err
			}
			for _, vv := range v.Values {
				_, err := dao.Value.Create(&model.Value{
					Value:       vv.Value.Value,
					Name:        vv.Name,
					Remark:      vv.Remark,
					DemandId:    vv.DemandId,
					Versioncode: vv.Versioncode,
					Imgs:        vv.Imgs,
					Operator:    op,
					CreatedTime: NowTime,
					UpdatedTime: NowTime,
					AttributeId: att_id,
					Status:      0,
				})
				if err != nil {
					return err
				}

			}
		}
	} else {
		////执行更新事件动作
		logger.Info("info", info.Status)
		if err := dao.Event.UpdateEvent(&model.Event{
			ID:           info.ID,
			Event:        info.Event.Event,
			Name:         info.Name,
			Versioncode:  info.Versioncode,
			Operator:     op,
			ReportDesc:   info.ReportDesc,
			Remark:       info.Remark,
			CreatedTime:  info.CreatedTime,
			UpdatedTime:  NowTime,
			Imgs:         info.Imgs,
			DemandId:     info.DemandId,
			Status:       info.Status,
			TempDemandId: info.TempDemandId,
		}); err != nil {
			return err
		}
		//更新分类
		//根据EventId删除原有的class-event记录

		if err := dao.Classification.DeleteByEventId(info.ID); err != nil {
			logger.Info(err.Error())
		}

		for _, v := range info.CategoryList {
			ClassEvent.Update(&model.ClassEvent{
				EventId:          info.Event.ID,
				ClassificationId: v,
			})
		}
		for _, v := range info.Properties {
			//如果属性id 为0 , 则为新增属性 ,包括新增 属性 取值
			if v.ID == 0 {
				//新增属性
				att_id, err := dao.Attribute.Create(&model.Attribute{
					Name:        v.Name,
					Key:         v.Key,
					Type:        v.Type,
					DemandId:    v.DemandId,
					VersionCode: v.VersionCode,
					Operator:    op,
					CreatedTime: NowTime,
					UpdatedTime: NowTime,
					Remark:      v.Remark,
					Imgs:        v.Imgs,
					EventId:     info.ID,
					Status:      0,
				})
				if err != nil {
					return err
				}
				//新增属性值
				for _, vv := range v.Values {
					_, err := dao.Value.Create(&model.Value{
						Value:       vv.Value.Value,
						Name:        vv.Name,
						Remark:      v.Remark,
						DemandId:    vv.DemandId,
						Versioncode: vv.Versioncode,
						Imgs:        v.Imgs,
						Operator:    op,
						CreatedTime: NowTime,
						UpdatedTime: NowTime,
						AttributeId: att_id,
						Status:      0,
					})
					if err != nil {
						return err
					}

				}
			} else {
				//更新属性
				if err := dao.Attribute.Update(&model.Attribute{
					ID:          v.ID,
					Name:        v.Name,
					Key:         v.Key,
					Type:        v.Type,
					DemandId:    v.DemandId,
					VersionCode: v.VersionCode,
					Operator:    op,
					CreatedTime: v.CreatedTime,
					UpdatedTime: NowTime,
					Remark:      v.Remark,
					Imgs:        v.Imgs,
					EventId:     info.ID,
					Status:      v.Status,
				}); err != nil {
					return err
				}
				//更新属性值
				for _, vv := range v.Values {
					if vv.ID == 0 {
						//新增属性值
						dao.Value.Create(&model.Value{
							Value:       vv.Value.Value,
							Name:        vv.Name,
							Remark:      vv.Remark,
							DemandId:    vv.DemandId,
							Versioncode: vv.Versioncode,
							Imgs:        vv.Imgs,
							Operator:    op,
							CreatedTime: NowTime,
							UpdatedTime: NowTime,
							AttributeId: v.ID, //属性id
							Status:      0,
						})
					} else {
						//更新属性值
						dao.Value.Update(&model.Value{
							ID:          vv.ID,
							Value:       vv.Value.Value,
							Name:        vv.Name,
							Remark:      vv.Remark,
							DemandId:    vv.DemandId,
							Versioncode: vv.Versioncode,
							Imgs:        vv.Imgs,
							Operator:    op,
							CreatedTime: vv.CreatedTime,
							UpdatedTime: NowTime,
							AttributeId: v.ID, //属性id
							Status:      vv.Status,
						})
					}
				}
			}
		}
	}
	return nil
}

//单个查询
func (*event) GetEventById(eventId uint) (*ResponseEventDetail, error) {
	//获取事件ID 所属的分类
	classIds, err := dao.Classevent.ClassIdByEventId(eventId)
	if err != nil {
		return nil, err
	}

	event, err := dao.Event.GetEventById(eventId)
	if err != nil {
		return nil, err
	}
	//定义一个属性list,用来存放属性和属性值
	attributelist := make([]*ResponseAttribute, 0)
	//根据eventid 获取attribute list
	attributeList, err := dao.Attribute.ListByEventId(eventId)
	if err != nil {
		return nil, err
	}

	//根据attributeid 获取value list
	for _, v := range attributeList {
		//根据属性id获取对应的value list
		values, err := dao.Value.ListByAttributeId(v.ID)
		if err != nil {
			return nil, err
		}
		//定义一个value list,用来存放value和tags
		valuelist := make([]*ResponseValue, 0)
		//遍历values list ,获取value信息 对应的tags
		for _, value := range values {
			v_tags, err := dao.Tags.ListByValueId(value.ID)
			if err != nil {
				return nil, err
			}
			v := ResponseValue{
				Value: model.Value{
					ID:          value.ID,
					Value:       value.Value,
					Name:        value.Name,
					Remark:      value.Remark,
					DemandId:    value.DemandId,
					Versioncode: value.Versioncode,
					Imgs:        value.Imgs,
					Operator:    value.Operator,
					CreatedTime: value.CreatedTime,
					UpdatedTime: value.UpdatedTime,
					AttributeId: value.AttributeId,
					Status:      value.Status,
				},
				Tags: v_tags,
			}
			valuelist = append(valuelist, &v)
		}

		//获取属性对应的tags
		a_tags, err := dao.Tags.ListByAttributeId(v.ID)
		//将value信息及tags 拼接到属性中
		attributes := &ResponseAttribute{
			Attribute: model.Attribute{
				ID:          v.ID,
				Name:        v.Name,
				Key:         v.Key,
				Type:        v.Type,
				DemandId:    v.DemandId,
				VersionCode: v.VersionCode,
				Operator:    v.Operator,
				CreatedTime: v.CreatedTime,
				UpdatedTime: v.UpdatedTime,
				Remark:      v.Remark,
				Imgs:        v.Imgs,
				EventId:     v.EventId,
				Status:      v.Status,
			},
			Values: valuelist,
			Tags:   a_tags,
		}
		attributelist = append(attributelist, attributes)
	}
	e_tags, err := dao.Tags.ListByEventId(eventId)
	if err != nil {
		return nil, err
	}
	info := &ResponseEventDetail{
		Event: model.Event{
			ID:          event.ID,
			Event:       event.Event,
			Name:        event.Name,
			Versioncode: event.Versioncode,
			ReportDesc:  event.ReportDesc,
			Remark:      event.Remark,
			CreatedTime: event.CreatedTime,
			UpdatedTime: event.UpdatedTime,
			Imgs:        event.Imgs,
			DemandId:    event.DemandId,
			Status:      event.Status,
		},
		Properties:   attributelist,
		CategoryList: classIds,
		Tags:         e_tags,
	}
	return info, nil
}

//并发查询
func (*event) GetEventById_Go(event *model.Event) (*ResponseEventDetail, error) {
	//获取事件ID 所属的分类
	classIds, err := dao.Classevent.ClassIdByEventId(event.ID)
	//logger.Info("获取事件ID 所属的分类", classIds)
	if err != nil {
		return nil, err
	}

	//定义一个属性list,用来存放属性和属性值
	attribute_list := make([]*ResponseAttribute, 0)
	//根据eventid 获取attribute list
	attributelist_data, err := dao.Attribute.ListByEventId(event.ID)
	//logger.Info("根据eventid 获取attribute list", len(attributelist_data))
	if err != nil {
		return nil, err
	}
	if len(attributelist_data) != 0 {
		var wg_0 sync.WaitGroup
		var lock_0 sync.Mutex
		for _, v := range attributelist_data {
			wg_0.Add(1)
			go func(v *model.Attribute) {
				defer wg_0.Done()

				//定义一个value list,用来存放value和tags
				valuelist := make([]*ResponseValue, 0)
				//根据属性id获取对应的value list
				values_data, err := dao.Value.ListByAttributeId(v.ID)
				//logger.Info("根据属性id获取对应的value list", len(values_data))
				if err != nil {
					return //nil, err
				}

				if len(values_data) != 0 {

					var wg sync.WaitGroup
					var lock sync.Mutex
					for _, value := range values_data {
						wg.Add(1)
						go func(value *model.Value) {
							defer wg.Done()
							v_tags, err := dao.Tags.ListByValueId(value.ID)
							//logger.Info("根据value id获取对应的tags", len(v_tags))
							if err != nil {
								return //nil, err
							}
							//升序排序
							sort.Slice(v_tags, func(i, j int) bool {
								if v_tags[i].ID < v_tags[j].ID {
									return true
								} else {
									return false
								}
							})
							v := ResponseValue{
								Value: model.Value{
									ID:          value.ID,
									Value:       value.Value,
									Name:        value.Name,
									Remark:      value.Remark,
									DemandId:    value.DemandId,
									Versioncode: value.Versioncode,
									Imgs:        value.Imgs,
									Operator:    value.Operator,
									CreatedTime: value.CreatedTime,
									UpdatedTime: value.UpdatedTime,
									AttributeId: value.AttributeId,
									Status:      value.Status,
								},
								Tags: v_tags,
							}
							lock.Lock()
							defer lock.Unlock()
							valuelist = append(valuelist, &v)
						}(value)

					}
					wg.Wait()
				}

				//获取属性对应的tags
				a_tags, err := dao.Tags.ListByAttributeId(v.ID)
				//logger.Info("获取属性id对应的tags", len(a_tags))
				//升序排序
				sort.Slice(valuelist, func(i, j int) bool {
					if valuelist[i].ID < valuelist[j].ID {
						return true
					} else {
						return false
					}
				})
				//升序排序
				sort.Slice(a_tags, func(i, j int) bool {
					if a_tags[i].ID < a_tags[j].ID {
						return true
					} else {
						return false
					}
				})
				//将value信息及tags 拼接到属性中
				attributes := &ResponseAttribute{
					Attribute: model.Attribute{
						ID:          v.ID,
						Name:        v.Name,
						Key:         v.Key,
						Type:        v.Type,
						DemandId:    v.DemandId,
						VersionCode: v.VersionCode,
						Operator:    v.Operator,
						CreatedTime: v.CreatedTime,
						UpdatedTime: v.UpdatedTime,
						Remark:      v.Remark,
						Imgs:        v.Imgs,
						EventId:     v.EventId,
						Status:      v.Status,
					},
					Values: valuelist,
					Tags:   a_tags,
				}
				lock_0.Lock()
				defer lock_0.Unlock()
				attribute_list = append(attribute_list, attributes)

			}(v)
		}
		wg_0.Wait()
	}
	//升序排序
	sort.Slice(attribute_list, func(i, j int) bool {
		if attribute_list[i].Attribute.ID < attribute_list[j].Attribute.ID {
			return true
		} else {
			return false
		}
	})
	e_tags, err := dao.Tags.ListByEventId(event.ID)
	if err != nil {
		return nil, err
	}
	info := &ResponseEventDetail{
		Event: model.Event{
			ID:           event.ID,
			Event:        event.Event,
			Name:         event.Name,
			Versioncode:  event.Versioncode,
			Operator:     event.Operator,
			ReportDesc:   event.ReportDesc,
			Remark:       event.Remark,
			CreatedTime:  event.CreatedTime,
			UpdatedTime:  event.UpdatedTime,
			Imgs:         event.Imgs,
			DemandId:     event.DemandId,
			Status:       event.Status,
			TempDemandId: event.TempDemandId,
			Kernel:       event.Kernel,
		},
		Properties:   attribute_list,
		CategoryList: classIds,
		Tags:         e_tags,
	}
	return info, nil
}

//根据搜索关键字获取event list
func (*event) ListBySearchKey(searchKey string) ([]*ResponseEventDetail, error) {
	//根据搜索关键字获取event list
	eventlist, err := dao.Event.EventIdsBySearchKeyword(searchKey)
	if err != nil {
		return nil, err
	}
	//根据搜索关键字获取属性对应的event list
	attributeEventList, err := dao.Attribute.EventIdsBySearchKeyword(searchKey)
	if err != nil {
		return nil, err
	}

	//根据搜索关键字获取属性对应的event list
	aIds, err := dao.Value.EventIdsBySearchKeyword(searchKey)
	if err != nil {
		return nil, err
	}
	//根据属性id获取eventid
	eids, err := dao.Attribute.GetEventidByaids(aIds)
	if err != nil {
		return nil, err
	}

	//将两个list合并
	eventlist = append(eventlist, attributeEventList...)
	eventlist = append(eventlist, eids...)
	//去重
	eventlist = RemoveDuplicateElement(eventlist)
	logger.Info(fmt.Sprintf("搜索关键字: %s,结果照片or属性: %d", searchKey, len(eventlist)))

	//定义一个事件list,用来存放事件和属性
	responseEventList := make([]*ResponseEventDetail, 0)

	//批量查询event信息,根据eventid
	events_data, err := dao.Event.ListByEventIds(eventlist)
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	var lock sync.Mutex
	//根据eventid 获取attribute list
	for _, event := range events_data {
		wg.Add(1)
		go func(event *model.Event) {
			defer wg.Done()
			info, err := Event.GetEventById_Go(event)
			if err != nil {
				return //nil, err
			}
			lock.Lock()
			defer lock.Unlock()
			responseEventList = append(responseEventList, info)
		}(event)

	}
	wg.Wait()
	//升序排序
	sort.Slice(responseEventList, func(i, j int) bool {
		if responseEventList[i].Event.ID < responseEventList[j].Event.ID {
			return true
		} else {
			return false
		}
	})
	return responseEventList, nil
}

//去重
func RemoveDuplicateElement(list []uint) []uint {
	result := make([]uint, 0, len(list))
	temp := map[uint]struct{}{}
	for _, item := range list {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

//根据versioncode获取event list
func (*event) ListByVersionCode(versioncode uint) ([]*ResponseEventDetail, error) {
	eventids, err := Event.ListAllIdByVersionCode(versioncode)

	//根据versioncode获取event list
	eventlist, err := dao.Event.ListByEventIds(eventids)
	if err != nil {
		return nil, err
	}

	//定义一个事件list,用来存放事件和属性
	responseEventList := make([]*ResponseEventDetail, 0)

	var wg sync.WaitGroup
	var lock sync.Mutex
	//遍历eventlist, 获取event信息
	for _, event := range eventlist {
		wg.Add(1)
		go func(event *model.Event) {
			defer wg.Done()
			info, err := Event.GetEventById_Go(event)
			if err != nil {
				return //nil, err
			}
			lock.Lock()
			defer lock.Unlock()
			responseEventList = append(responseEventList, info)
		}(event)

	}
	wg.Wait()
	//升序排序
	sort.Slice(responseEventList, func(i, j int) bool {
		if responseEventList[i].Event.ID < responseEventList[j].Event.ID {
			return true
		} else {
			return false
		}
	})
	return responseEventList, nil
}

//根据versioncode 获取所有的event ids
func (*event) ListAllIdByVersionCode(versioncode uint) ([]uint, error) {

	eventid_list := make([]uint, 0)
	//tag中包括versioncode,获取所有的attid_list
	tagslist, err := dao.Tags.ListByVersionCode(versioncode)
	if err != nil {
		return nil, err
	}
	attid_list := make([]uint, 0)
	valueid_list := make([]uint, 0)
	for _, v := range tagslist {
		if v.AttributeId != 0 {
			attid_list = append(attid_list, v.AttributeId)
		}
		if v.ValueId != 0 {
			valueid_list = append(valueid_list, v.ValueId)
		}
	}
	//--根据valueid_list获取attid_list
	valuedata_list, err := dao.Value.ListByValueIds(valueid_list)
	if err != nil {
		return nil, err
	}
	for _, v := range valuedata_list {
		attid_list = append(attid_list, v.AttributeId)
	}
	//logger.Info("tags 获取attid_list", len(attid_list))
	//取值中包括versioncode,获取所有的attid_list
	valuedata_list, err = dao.Value.ListByVersionCode(versioncode)
	if err != nil {
		return nil, err
	}
	for _, v := range valuedata_list {
		attid_list = append(attid_list, v.AttributeId)
	}
	//logger.Info("value 获取attid_list", len(attid_list))
	//属性中包括versioncode,获取所有的attid_list
	attributedata_list, err := dao.Attribute.ListByVersionCode(versioncode)
	if err != nil {
		return nil, err
	}
	for _, v := range attributedata_list {
		attid_list = append(attid_list, v.ID)
	}
	//logger.Info("att 获取attid_list", len(attid_list))
	attdata_list, err := dao.Attribute.ListByAttributeIds(attid_list)
	if err != nil {
		return nil, err
	}
	//获取所有的eventid_list
	for _, v := range attdata_list {
		eventid_list = append(eventid_list, v.EventId)
	}
	//logger.Info("获取tags/属性/取值涉及到的eventid_list", len(eventid_list))
	//获取event中包括versioncode
	eventdata_list, err := dao.Event.ListByVersionCode(versioncode)
	if err != nil {
		return nil, err
	}
	for _, v := range eventdata_list {
		eventid_list = append(eventid_list, v.ID)
	}
	//logger.Info("获取所有eventid_list", len(eventid_list))
	//去重
	eventid_list = RemoveDuplicateElement(eventid_list)
	//logger.Info("获取所有eventid_list,去重", len(eventid_list))

	return eventid_list, nil
}

//根据demandid 获取event list
func (*event) ListByDemandId(DemandId uint) ([]*ResponseEventDetail, error) {
	eventids_list, err := Event.ListAllIdByDemandId(DemandId)
	if err != nil {
		return nil, err
	}
	//根据eventids_list 获取event list
	eventlist, err := dao.Event.ListByEventIds(eventids_list)
	if err != nil {
		return nil, err
	}
	//定义一个事件list,用来存放事件和属性
	responseEventList := make([]*ResponseEventDetail, 0)

	//根据eventid 获取attribute list
	var wg sync.WaitGroup
	var lock sync.Mutex
	for _, event := range eventlist {
		wg.Add(1)

		go func(event *model.Event) {
			defer wg.Done()
			info, err := Event.GetEventById_Go(event)
			if err != nil {
				return //nil, err
			}
			lock.Lock()
			responseEventList = append(responseEventList, info)
			defer lock.Unlock()
		}(event)

	}
	wg.Wait()
	//升序排序
	sort.Slice(responseEventList, func(i, j int) bool {
		if responseEventList[i].Event.ID < responseEventList[j].Event.ID {
			return true
		} else {
			return false
		}
	})
	return responseEventList, nil
}

//根据demandid 获取所有的event ids
func (*event) ListAllIdByDemandId(DemandId uint) ([]uint, error) {
	eventid_list := make([]uint, 0)

	//value 中包括demandid,获取所有的attid_list
	valuedata_list, err := dao.Value.ListByDemandId(DemandId)
	if err != nil {
		return nil, err
	}
	attid_list := make([]uint, 0)
	for _, v := range valuedata_list {
		attid_list = append(attid_list, v.AttributeId)
	}
	//attribute 中包括demandid,获取所有的attid_list
	attributedata_list, err := dao.Attribute.ListByDemandId(DemandId)
	if err != nil {
		return nil, err
	}
	for _, v := range attributedata_list {
		attid_list = append(attid_list, v.ID)
	}

	//根据attid_list获取eventid_list
	attributedata_list, err = dao.Attribute.ListByAttributeIds(attid_list)
	if err != nil {
		return nil, err
	}
	for _, v := range attributedata_list {
		eventid_list = append(eventid_list, v.EventId)
	}
	//事件中包括demandid,获取所有的eventid_list
	eventdata_list, err := dao.Event.ListByDemandId(DemandId)
	if err != nil {
		return nil, err
	}
	for _, v := range eventdata_list {
		eventid_list = append(eventid_list, v.ID)
	}
	//去重
	eventid_list = RemoveDuplicateElement(eventid_list)

	return eventid_list, nil

}

//RealDelete
func (*event) RealDelete(eventid uint) error {
	//删除event
	err := dao.Event.Delete(&model.Event{ID: eventid})
	if err != nil {
		return err
	}
	//删除分类
	err = dao.Classevent.DeleteByEventId(eventid)
	if err != nil {
		return err
	}
	//查询出attid
	attdata_list, err := dao.Attribute.ListByEventId(eventid)
	if err != nil {
		return err
	}
	attid_list := make([]uint, 0)
	for _, v := range attdata_list {
		attid_list = append(attid_list, v.ID)
	}
	//删除属性
	err = dao.Attribute.DeleteByIds(attid_list)
	if err != nil {
		return err
	}
	//根据attid_list查询出valueid_list
	valuedata_list, err := dao.Value.ListByAttributeIds(attid_list)
	valueid_list := make([]uint, 0)
	for _, v := range valuedata_list {
		valueid_list = append(valueid_list, v.ID)
	}
	//删除属性值
	err = dao.Value.DeleteByIds(valueid_list)
	if err != nil {
		return err
	}
	//删除tag
	err = dao.Tags.DeleteByAttributeId(attid_list)
	if err != nil {
		return err
	}
	err = dao.Tags.DeleteByValueId(valueid_list)
	if err != nil {
		return err
	}
	err = dao.Tags.DeleteEventId(eventid)
	if err != nil {
		return err
	}
	return nil
}

//Remove
func (*event) RemoveDemand(eventId, demandId uint) error {
	//查询出event
	err := dao.Event.UpdateEventDemand(eventId)
	if err != nil {
		return err
	}

	//更新属性
	err = dao.Attribute.UpdateDemand(eventId, demandId)
	if err != nil {
		return err
	}
	//更新属性
	err = dao.Value.UpdateDemand(eventId, demandId)
	if err != nil {
		return err
	}

	return nil
}

func (*event) ListKernel() ([]*ResponseEventDetail, error) {

	eventlist, err := dao.Event.ListKernel()
	if err != nil {
		return nil, err
	}
	//定义一个事件list,用来存放事件和属性
	responseEventList := make([]*ResponseEventDetail, 0)

	//根据eventid 获取attribute list
	var wg sync.WaitGroup
	var lock sync.Mutex
	for _, event := range eventlist {
		wg.Add(1)

		go func(event *model.Event) {
			defer wg.Done()
			info, err := Event.GetEventById_Go(event)
			if err != nil {
				return //nil, err
			}
			lock.Lock()
			responseEventList = append(responseEventList, info)
			defer lock.Unlock()
		}(event)

	}
	wg.Wait()
	//升序排序
	sort.Slice(responseEventList, func(i, j int) bool {
		if responseEventList[i].Event.ID < responseEventList[j].Event.ID {
			return true
		} else {
			return false
		}
	})
	return responseEventList, nil
}
