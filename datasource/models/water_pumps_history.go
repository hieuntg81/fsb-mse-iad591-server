package models

import (
	"gorm.io/gorm"
	"time"
)

func (WaterPumpsHistory) TableName() string {
	return "water_pumps_histories"
}

type WaterPumpsHistory struct {
	gorm.Model
	ID        uint `gorm:"primaryKey;autoIncrement:false;column:id"`
	OpenTime  uint `gorm:"column:open_time"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UpdatedBy string
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
