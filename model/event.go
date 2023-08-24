package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

//事件表   是和需求走还是和版本走?
type Event struct {
	ID uint `json:"eventId" gorm:"primary_key"`
	//CreatedAt time.Time `json:"createdTime"`
	//UpdatedAt time.Time `json:"updatedTime"`
	//DeletedAt *time.Time `sql:"index"`

	Event       string `json:"event" gorm:"column:event"`             //事件 英文名称
	Name        string `json:"name" gorm:"column:name"`               //事件 中文名称
	Versioncode uint   `json:"versioncode" gorm:"column:versioncode"` //版本号
	Operator    string `json:"operator" gorm:"column:operator"`       //操作人
	ReportDesc  string `json:"reportDesc" gorm:"column:reportDesc"`   //上报时机
	Remark      string `json:"remark" gorm:"column:remark"`           //事件备注
	CreatedTime int64  `json:"createdTime" gorm:"column:createdTime"` //创建时间
	UpdatedTime int64  `json:"updatedTime" gorm:"column:updatedTime"` //更新时间
	Imgs        Array  `json:"imgs" gorm:"column:imgs" `              //配图

	DemandId uint `json:"demandId" gorm:"column:demandId"` //事件对应的需求id,需求在demand表中

	Status       uint `json:"status" gorm:"column:status"`             //事件状态 0:正常 1:已删除
	TempDemandId uint `json:"tempDemandId" gorm:"column:tempDemandId"` //临时需求id, 用于添加事件到其他需求时, 临时存储需求id
	Kernel       bool `json:"kernel" gorm:"column:kernel"`             //是否是核心事件

	//点击创建时间时, 创建对应的事件记录, 另创建 classevent 表中 对应的分类id --> 事件id  , 一个事件可属于多个分类
}

func (*Event) TableName() string {
	return "event"
}

type Array []string

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (a *Array) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to scan Array value:", value))
	}
	*a = strings.Split(string(bytes), ",")
	return nil
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (a Array) Value() (driver.Value, error) {
	if len(a) > 0 {
		var str string = a[0]
		for _, v := range a[1:] {
			str += "," + v
		}
		return str, nil
	} else {
		return "", nil
	}
}
