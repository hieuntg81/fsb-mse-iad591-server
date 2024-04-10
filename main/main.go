package main

import (
	"fmt"
	"fsb-mse-iad591-server/datasource"
	"fsb-mse-iad591-server/datasource/models"
	controller "fsb-mse-iad591-server/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	DB := datasource.Init()
	err := DB.AutoMigrate(&models.TemperatureRecord{}, &models.HumidityRecord{}, &models.WaterPumpsHistory{})
	if err != nil {
		fmt.Println(err)
		return
	}

	consumer := datasource.InitKafkaConsumer()
	err = consumer.SubscribeTopics([]string{"test_update_humidity"}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	producer := datasource.InitKafkaProducer()
	kh := controller.NewKafkaHandler(producer, consumer)

	h := controller.New(DB)
	router := gin.Default()
	api := router.Group("")
	{
		api.GET("/temperature", h.GetTemperatureRecords)
		api.GET("/humidity", h.GetHumidityRecords)
		api.GET("/humidity/test", kh.TestAddHumidity)
	}
	err = router.Run(":8080")
	if err != nil {
		return
	}
}
