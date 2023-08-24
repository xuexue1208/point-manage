package model

import "time"

//分类表
type ClassEvent struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	ClassificationId uint `json:"classification_id"  gorm:"column:classification_id"` //分类id, classification表中的id
	EventId          uint `json:"event_id"  gorm:"column:event_id"`                   //事件id, event表中的id

}

func (*ClassEvent) TableName() string {
	return "classevent"
}
