package datasource

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// 2181|2888|3888
func InitKafkaProducer() *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "138.2.53.115:2181"})
	if err != nil {
		panic(err)
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	return p
}

func InitKafkaConsumer() *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "138.2.53.115:2181",
		"group.id":          "default_group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	return c
}
