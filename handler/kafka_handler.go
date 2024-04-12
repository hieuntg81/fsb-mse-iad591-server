package handler

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"gorm.io/gorm"
	"iad591/datasource/models"
	"log"
	"time"
)

type KafkaHandler struct {
	producer sarama.SyncProducer
}

func NewKafkaClient(producer sarama.SyncProducer) KafkaHandler {
	return KafkaHandler{producer}
}

func (kh KafkaHandler) SendLatestHumidity(DB *gorm.DB) {
	var records []models.TemperatureRecord
	if result := DB.Find(&records); result.Error != nil {
		fmt.Println(result.Error)
	}
	firstElement := records[0]
	data, err := json.Marshal(firstElement)
	msg := &sarama.ProducerMessage{
		Topic:     "bizfly-7-678-configuration",
		Value:     sarama.StringEncoder(data),
		Timestamp: time.Now(),
	}

	par, offset, err := kh.producer.SendMessage(msg)
	fmt.Printf("Partition %d", &par)
	fmt.Printf("Offset %d", &offset)
	if err != nil {
		log.Fatalln("Failed to send message:", err)
	}
}
