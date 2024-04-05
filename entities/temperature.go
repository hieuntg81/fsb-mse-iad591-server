package entities

import (
	"github.com/google/uuid"
	"time"
)

type TemperatureRecord struct {
	ID        string    `db:"id"`
	Value     string    `db:"value"`
	Unit      string    `db:"unit"`
	CreatedAt time.Time `db:"created_date"`
	UpdatedAt time.Time `db:"last_modified_date"`
	UpdatedBy string    `db:"last_modified_by"`
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
