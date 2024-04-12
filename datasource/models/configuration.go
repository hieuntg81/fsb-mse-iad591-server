package models

import (
	"gorm.io/gorm"
	"time"
)

func (Configuration) TableName() string {
	return "configuration"
}

type Configuration struct {
	gorm.Model
	ID              uint `gorm:"primaryKey;autoIncrement:true;column:id"`
	Type            uint
	Time            uint
	CurrentHumidity float64
	TriggerHumidity float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	UpdatedBy       string
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type ConfigurationPayload struct {
	Type            uint    `json:"type"`
	Time            uint    `json:"time"`
	CurrentHumidity float64 `json:"current_humidity"`
	TriggerHumidity float64 `json:"trigger_humidity"`
	UpdateBy        string  `json:"updated_by"`
}
