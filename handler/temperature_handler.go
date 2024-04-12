package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	models2 "iad591/datasource/models"
	"io"
	"net/http"
	"time"
)

type Handler struct {
	DB      *gorm.DB
	handler KafkaHandler
}

func New(db *gorm.DB, handler KafkaHandler) Handler {
	return Handler{db, handler}
}

func (h Handler) GetTemperatureRecords(ctx *gin.Context) {
	var records []models2.TemperatureRecord
	if result := h.DB.Find(&records); result.Error != nil {
		fmt.Println(result.Error)
	}
	ctx.IndentedJSON(http.StatusOK, records)
}

func (h Handler) UpdateTemperatureRecords(ctx *gin.Context) {
	var payload models2.Payload
	asa, _ := io.ReadAll(ctx.Request.Body)
	err := json.Unmarshal(asa, &payload)
	if err != nil {
		return
	}
	var record models2.TemperatureRecord
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()
	record.Unit = payload.Unit
	record.Value = payload.Value
	record.UpdatedBy = payload.UpdatedBy
	if result := h.DB.Create(&record); result.Error != nil {
		fmt.Println(result.Error)
	}
}
