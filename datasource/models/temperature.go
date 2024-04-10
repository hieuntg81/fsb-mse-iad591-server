package models

import (
	"gorm.io/gorm"
	"time"
)

func (TemperatureRecord) TableName() string {
	return "temperatures"
}

type TemperatureRecord struct {
	gorm.Model
	ID        uint `gorm:"primaryKey;autoIncrement:false;column:id"`
	Value     string
	Unit      string
	CreatedAt time.Time
	UpdatedAt time.Time
	UpdatedBy string
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
