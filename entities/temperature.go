package entities

import (
	"github.com/google/uuid"
	"time"
)

type TemperatureRecord struct {
	ID        string
	Value     string
	Unit      string
	CreatedAt time.Time
	UpdatedAt time.Time
	UpdatedBy string
}

func NewTemperatureRecord(value, unit, updatedBy string) *TemperatureRecord {
	return &TemperatureRecord{
		ID:        uuid.New().String(),
		Unit:      unit,
		Value:     value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UpdatedBy: updatedBy,
	}
}
