package model

import "gorm.io/gorm"

type AdEntity struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	RedirectUrl string `gorm:"column:redirect_url"`
	Text        string `gorm:"column:text"`
}

// TableName 指定表名，对应Java的@Table注解
func (AdEntity) TableName() string {
	return "ads"
}
