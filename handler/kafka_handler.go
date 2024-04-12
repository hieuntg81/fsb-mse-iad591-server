package handler

import (
	"context"
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
	consumer sarama.ConsumerGroup
	DB       *gorm.DB
}

func NewKafkaClient(producer sarama.SyncProducer, consumer sarama.ConsumerGroup, db *gorm.DB) KafkaHandler {
	return KafkaHandler{producer, consumer, db}
}

func (kh KafkaHandler) SendLatestConfiguration(DB *gorm.DB) {
	var records []models.Configuration
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

func (kh KafkaHandler) UpdateHumidity(DB *gorm.DB) {
	var topic = "bizfly-7-678-humidity"
	consumer := KafkaConsumerGroupHandler{
		DB: kh.DB,
	}
	ctx := context.Background()

	for {
		err := kh.consumer.Consume(ctx, []string{topic}, consumer)
		if err != nil {
			log.Panicf("Error from consumer: %v", err)
		}
	}
}

// test only
func (kh KafkaHandler) TestSendHumidity(DB *gorm.DB) {
	var records []models.HumidityRecord
	if result := DB.Find(&records); result.Error != nil {
		fmt.Println(result.Error)
	}
	firstElement := records[0]
	data, err := json.Marshal(firstElement)
	msg := &sarama.ProducerMessage{
		Topic:     "bizfly-7-678-humidity",
		Value:     sarama.StringEncoder(data),
		Timestamp: time.Now(),
	}

	par, offset, err := kh.producer.SendMessage(msg)
	fmt.Printf("Partition %d ", &par)
	fmt.Printf("Offset %d", &offset)
	if err != nil {
		log.Fatalln("Failed to send message:", err)
	}
}

type KafkaConsumerGroupHandler struct {
	DB *gorm.DB
}

func (kcg KafkaConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	// Mark the beginning of a new session. This is called when the consumer group is being rebalanced.
	log.Println("Consumer group is being rebalanced")
	return nil
}

func (kcg KafkaConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	// Mark the end of the current session. This is called just before the next rebalance happens.
	log.Println("Rebalancing will happen soon, current session will end")
	return nil
}

func (kcg KafkaConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// This is where you put your message handling logic
	for message := range claim.Messages() {
		var payload models.Payload
		log.Printf("Message topic:%q partition:%d offset:%d\n", message.Topic, message.Partition, message.Offset)
		log.Printf("Received message from topic:%s partition:%d offset:%d key:%v value:%s", message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))
		err := json.Unmarshal(message.Value, &payload)
		if err != nil {
			log.Fatalln("Failed to parse message:", err)
		}
		var record models.HumidityRecord
		record.CreatedAt = time.Now()
		record.UpdatedAt = time.Now()
		record.Unit = payload.Unit
		record.Value = payload.Value
		record.UpdatedBy = "kafka"
		if result := kcg.DB.Create(&record); result.Error != nil {
			fmt.Println(result.Error)
		}
		sess.MarkMessage(message, "")
	}
	return nil
}
