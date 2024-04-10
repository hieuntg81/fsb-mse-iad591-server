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
	ID        uint `gorm:"primaryKey;autoIncrement:true;column:id"`
	Value     float64
	Unit      string
	CreatedAt time.Time
	UpdatedAt time.Time
	UpdatedBy string
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
