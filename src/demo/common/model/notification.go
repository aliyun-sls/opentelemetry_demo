package model

import "gorm.io/gorm"

type NotificationEntity struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	RedirectUrl string `gorm:"column:redirect_url"`
	Text        string `gorm:"column:text"`
}

// TableName 指定表名，对应Java的@Table注解
func (NotificationEntity) TableName() string {
	return "notification"
}
