package service

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestConsumerkafka(t *testing.T) {
	cfg := loadJsonConfig()
	consumer := DoInitConsumer(cfg)

	consumer.SubscribeTopics([]string{cfg.Topic}, nil)

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			var logisticsMsg = &LogisticsMsg{}
			err = json.Unmarshal(msg.Value, logisticsMsg)
			if err != nil {
				fmt.Printf("Consumer Unmarshal error: %v (%v)\n", err, msg)
			}
			marshal, _ := json.Marshal(logisticsMsg)
			fmt.Println(string(marshal))
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	consumer.Close()
}
