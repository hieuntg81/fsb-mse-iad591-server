package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/xdg/scram"
	"iad591/datasource"
	"iad591/datasource/models"
	controller "iad591/handler"
	"log"
	"time"
)

func disableCORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")  // You can set specific origins here (e.g., "http://localhost:3000")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "*") // Allow all methods
	c.Writer.Header().Set("Access-Control-Allow-Headers", "*") // Allow all headers
	c.Next()
}

func main() {

	DB := datasource.Init()
	err := DB.AutoMigrate(&models.TemperatureRecord{}, &models.Configuration{}, &models.HumidityRecord{}, &models.WaterPumpsHistory{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create producer
	kh := controller.NewKafkaClient(InitKafkaProducer(), InitKafkaConsumer(), DB)
	schedule := cron.New()
	_, err = schedule.AddFunc("@every 1m", func() {
		kh.SendLatestConfiguration(DB)
		//kh.TestSendHumidity(DB)
		fmt.Println("Sending configuration every minutes")
	})
	go schedule.Start()

	// Create consumer
	go kh.UpdateHumidity(DB)

	h := controller.New(DB, kh)
	router := gin.Default()
	router.Use(disableCORS)
	api := router.Group("")
	{
		api.GET("/temperature", h.GetTemperatureRecords)
		api.GET("/humidity", h.GetHumidityRecords)
		api.POST("/humidity", h.UpdateHumidityRecords)
		api.POST("/temperature", h.UpdateTemperatureRecords)
		api.GET("/waterpumping", h.GetWaterPumpsHistories)
		api.POST("/waterpumping", h.UpdateWaterPumpsHistories)
		api.GET("/configuration", h.GetConfiguration)
		api.POST("/configuration", h.UpdateConfiguration)
	}
	err = router.Run(":6666")
	if err != nil {
		return
	}
}

func InitKafkaProducer() sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Net.SASL.Enable = true
	config.Net.SASL.User = "bizfly-7-678-admin"
	config.Net.SASL.Password = "admin1234"
	config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
	config.Net.SASL.Handshake = true
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_0_0_0

	brokers := []string{"cluster001.kas.bfcplatform.vn:9094"}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}
	return producer
}

func InitKafkaConsumer() sarama.ConsumerGroup {
	config := sarama.NewConfig()
	config.Net.SASL.Enable = true
	config.Net.SASL.User = "bizfly-7-678-admin"
	config.Net.SASL.Password = "admin1234"
	config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
	config.Net.SASL.Handshake = true
	config.Version = sarama.V2_0_0_0
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	brokers := []string{"cluster001.kas.bfcplatform.vn:9094"}
	//consumer, err := sarama.NewConsumer(brokers, config)
	//if err != nil {
	//	log.Fatalln("Error creating consumer", err)
	//}
	consumerGroup, err := sarama.NewConsumerGroup(brokers, "bizfly-7-678-group-1", config)
	if err != nil {
		log.Panicf("Error creating consumerGroup group client: %v", err)
	}
	return consumerGroup
}

var (
	SHA256 scram.HashGeneratorFcn = sha256.New
	SHA512 scram.HashGeneratorFcn = sha512.New
)

type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}
