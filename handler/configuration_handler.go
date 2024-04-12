package handler

import (
	"encoding/json"
	"fmt"
	models2 "fsb-mse-iad591-server/datasource/models"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

func (h Handler) GetConfiguration(ctx *gin.Context) {
	var records []models2.Configuration
	if result := h.DB.Find(&records); result.Error != nil {
		fmt.Println(result.Error)
	}
	ctx.IndentedJSON(http.StatusOK, records)
}

func (h Handler) UpdateConfiguration(ctx *gin.Context) {
	h.DB.Exec("DELETE FROM configuration WHERE true")
	var payload models2.ConfigurationPayload
	asa, _ := io.ReadAll(ctx.Request.Body)
	err := json.Unmarshal(asa, &payload)
	if err != nil {
		return
	}
	var record models2.Configuration
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()
	record.Type = payload.Type
	record.Time = payload.Time
	record.Time = payload.Time
	record.CurrentHumidity = payload.CurrentHumidity
	record.TriggerHumidity = payload.TriggerHumidity
	record.UpdatedBy = payload.UpdateBy
	if result := h.DB.Create(&record); result.Error != nil {
		fmt.Println(result.Error)
	}
}
