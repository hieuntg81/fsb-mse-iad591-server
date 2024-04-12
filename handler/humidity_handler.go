package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"iad591/datasource/models"
	"io"
	"net/http"
	"time"
)

func (h Handler) GetHumidityRecords(ctx *gin.Context) {
	var records []models.HumidityRecord
	if result := h.DB.Find(&records); result.Error != nil {
		fmt.Println(result.Error)
	}
	ctx.IndentedJSON(http.StatusOK, records)
}

func (h Handler) UpdateHumidityRecords(ctx *gin.Context) {
	var payload models.Payload
	asa, _ := io.ReadAll(ctx.Request.Body)
	err := json.Unmarshal(asa, &payload)
	if err != nil {
		return
	}
	var record models.HumidityRecord
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()
	record.Unit = payload.Unit
	record.Value = payload.Value
	record.UpdatedBy = payload.UpdatedBy
	if result := h.DB.Create(&record); result.Error != nil {
		fmt.Println(result.Error)
	}
}
