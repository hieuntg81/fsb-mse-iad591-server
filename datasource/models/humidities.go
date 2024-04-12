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
	ID        uint `gorm:"primaryKey;autoIncrement:true;column:id"`
	Value     float64
	Unit      string
	CreatedAt time.Time
	UpdatedAt time.Time
	UpdatedBy string
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Payload struct {
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
	UpdatedBy string  `json:"updated_by"`
}

type PumpPayload struct {
	CurrentHumidity float64 `json:"current_humidity"`
	OpenTime        uint    `json:"open_time"`
	UpdatedBy       string  `json:"updated_by"`
}
