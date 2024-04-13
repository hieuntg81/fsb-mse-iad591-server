package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	models "iad591/datasource/models"
	"io"
	"net/http"
	"time"
)

func (h Handler) GetConfiguration(ctx *gin.Context) {
	var records []models.Configuration
	if result := h.DB.Find(&records); result.Error != nil {
		fmt.Println(result.Error)
	}
	ctx.IndentedJSON(http.StatusOK, records)
}

func (h Handler) UpdateConfiguration(ctx *gin.Context) {
	var records []models.Configuration
	if result := h.DB.Find(&records); result.Error != nil {
		fmt.Println(result.Error)
	}
	record := records[0]
	var payload models.ConfigurationPayload
	asa, _ := io.ReadAll(ctx.Request.Body)
	err := json.Unmarshal(asa, &payload)
	if err != nil {
		return
	}

	record.UpdatedAt = time.Now()
	if payload.Type != 0 {
		record.Type = payload.Type
	}

	if payload.Time != 0 {
		record.Time = payload.Time
	}

	if payload.CurrentHumidity != 0 {
		record.CurrentHumidity = payload.CurrentHumidity
	}

	if payload.TriggerHumidity != 0 {
		record.TriggerHumidity = payload.TriggerHumidity
	}
	record.UpdatedBy = payload.UpdateBy
	if result := h.DB.Updates(&record); result.Error != nil {
		fmt.Println(result.Error)
	}
	ctx.IndentedJSON(http.StatusOK, record)
}
