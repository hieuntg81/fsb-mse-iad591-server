package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"fsb-mse-iad591-server/datasource"
	models2 "fsb-mse-iad591-server/datasource/models"
	controller "fsb-mse-iad591-server/handler"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/xdg-go/scram"
	"log"
	"os"
	"os/signal"
)

func main() {
	DB := datasource.Init()
	err := DB.AutoMigrate(&models2.TemperatureRecord{}, &models2.Configuration{}, &models2.HumidityRecord{}, &models2.WaterPumpsHistory{})
	if err != nil {
		fmt.Println(err)
		return
	}

	//initKafkaConsumer()

	h := controller.New(DB)
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
	if err != nil {
		return
	}
}

func initKafkaConsumer() {
	conf := sarama.NewConfig()
	conf.Net.SASL.Enable = true
	conf.Net.SASL.User = "bizfly-7-678-admin" // Thay thế bằng tài khoản của bạn
	conf.Net.SASL.Password = "admin1234"      // Thay thế bằng mật khẩu của bạn
	conf.Net.SASL.Handshake = true

	// Sử dụng SCRAM-SHA-512
	conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
		return &XDGSCRAMClient{HashGeneratorFcn: SHA512}
	}
	conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	brokers := []string{"cluster001.kas.bfcplatform.vn:9094"}
	consumer, err := sarama.NewConsumer(brokers, conf)

	if err != nil {
		panic(err)
	}

	// Tạo một channel để nhận tín hiệu ngắt
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Consume từ topic "my_topic"
	partition, err := consumer.ConsumePartition("bizfly-7-678-mse-topic", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	// Đọc message
	go func() {
		for {
			select {
			case err := <-partition.Errors():
				log.Println(err)
			case msg := <-partition.Messages():
				log.Printf("Received message: %s\n", string(msg.Value))
			}
		}
	}()

	// Đợi tín hiệu ngắt
	<-signals
	err = consumer.Close()
	if err != nil {
		return
	}
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
