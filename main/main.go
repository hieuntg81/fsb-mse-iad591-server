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
)

func main() {
	DB := datasource.Init()
	err := DB.AutoMigrate(&models.TemperatureRecord{}, &models.Configuration{}, &models.HumidityRecord{}, &models.WaterPumpsHistory{})
	if err != nil {
		fmt.Println(err)
		return
	}

	h := controller.New(DB)
	kh := controller.NewKafkaClient(InitKafkaProducer())
	router := gin.Default()
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

	scheduler := cron.New()
	_, err = scheduler.AddFunc("0 * * * * ?", func() {
		kh.SendLatestHumidity(DB)
		fmt.Println("Sending configuration every minutes")
	})
	scheduler.Run()
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
