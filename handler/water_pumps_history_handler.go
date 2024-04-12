package handler

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	models2 "iad591/datasource/models"
	"io"
	"log"
	"net/http"
	"time"
)

func (h Handler) GetWaterPumpsHistories(ctx *gin.Context) {
	var records []models2.WaterPumpsHistory
	if result := h.DB.Find(&records); result.Error != nil {
		fmt.Println(result.Error)
	}
	ctx.IndentedJSON(http.StatusOK, records)
}

func (h Handler) UpdateWaterPumpsHistories(ctx *gin.Context) {
	var payload models2.PumpPayload
	asa, _ := io.ReadAll(ctx.Request.Body)
	err := json.Unmarshal(asa, &payload)
	if err != nil {
		return
	}
	var record models2.WaterPumpsHistory
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()
	record.OpenTime = payload.OpenTime
	record.CurrentHumidity = payload.CurrentHumidity
	record.UpdatedBy = payload.UpdatedBy
	if result := h.DB.Create(&record); result.Error != nil {
		fmt.Println(result.Error)
	}

	data, err := json.Marshal(record)
	msg := &sarama.ProducerMessage{
		Topic:     "bizfly-7-678-water-pumps-trigger",
		Value:     sarama.StringEncoder(data),
		Timestamp: time.Now(),
	}

	par, offset, err := h.handler.producer.SendMessage(msg)
	fmt.Printf("Partition %d ", &par)
	fmt.Printf("Offset %d", &offset)
	if err != nil {
		log.Fatalln("Failed to send message:", err)
	}
}
