package handler

import (
	"encoding/json"
	"fmt"
	"fsb-mse-iad591-server/datasource/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Handler) GetHumidityRecords(ctx *gin.Context) {
	var records []models.HumidityRecord
	if result := h.DB.Find(&records); result.Error != nil {
		fmt.Println(result.Error)
	}
	ctx.IndentedJSON(http.StatusOK, records)
}

func (h Handler) UpdateHumidityRecords(body string) {
	var record models.HumidityRecord
	err := json.Unmarshal([]byte(body), &record)
	if err != nil {
		return
	}

	if result := h.DB.Create(&record); result.Error != nil {
		fmt.Println(result.Error)
	}
}
