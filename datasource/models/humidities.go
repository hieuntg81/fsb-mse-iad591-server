package models

import (
	"gorm.io/gorm"
	"time"
)

func (HumidityRecord) TableName() string {
	return "humidities"
}

type HumidityRecord struct {
	gorm.Model
	ID        uint `gorm:"primaryKey;autoIncrement:false;column:id"`
	Value     string
	Unit      string
	CreatedAt time.Time
	UpdatedAt time.Time
	UpdatedBy string
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
