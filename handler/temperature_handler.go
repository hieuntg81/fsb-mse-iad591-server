package handler

import (
	"encoding/json"
	"fmt"
	"fsb-mse-iad591-server/datasource/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
	"time"
)

type Handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Handler {
	return Handler{db}
}

func (h Handler) GetTemperatureRecords(ctx *gin.Context) {
	var records []models.TemperatureRecord
	if result := h.DB.Find(&records); result.Error != nil {
		fmt.Println(result.Error)
	}
	ctx.IndentedJSON(http.StatusOK, records)
}

func (h Handler) UpdateTemperatureRecords(ctx *gin.Context) {
	var payload models.Payload
	asa, _ := io.ReadAll(ctx.Request.Body)
	err := json.Unmarshal(asa, &payload)
	if err != nil {
		return
	}
	var record models.TemperatureRecord
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()
	record.Unit = payload.Unit
	record.Value = payload.Value
	record.UpdatedBy = payload.UpdatedBy
	if result := h.DB.Create(&record); result.Error != nil {
		fmt.Println(result.Error)
	}
}
