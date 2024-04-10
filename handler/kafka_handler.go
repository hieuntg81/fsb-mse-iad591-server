package handler

import (
	"fmt"
	"fsb-mse-iad591-server/datasource/models"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type KafkaHandler struct {
	producer *kafka.Producer
	consumer *kafka.Consumer
}

func NewKafkaHandler(p *kafka.Producer, c *kafka.Consumer) KafkaHandler {
	return KafkaHandler{p, c}
}

func (kh KafkaHandler) TestAddHumidity(ctx *gin.Context) {
	record := models.HumidityRecord{
		Value:     "1",
		Unit:      "%",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UpdatedBy: "Kafka",
		DeletedAt: gorm.DeletedAt{},
	}

	topic := "test_update_humidity"

	err := kh.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(fmt.Sprintf("%v", record)),
	}, nil)

	if err != nil {
		fmt.Print(err)
		return
	}
	ctx.String(http.StatusOK, "Done")
}
