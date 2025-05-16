package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"sls-mall-go/common/util"
)

type KafkaConfig struct {
	Topic            string `json:"topic"`
	GroupId          string `json:"group.id"`
	BootstrapServers string `json:"bootstrap.servers"`
	SecurityProtocol string `json:"security.protocol"`
	SslCaLocation    string `json:"ssl.ca.location"`
	SaslMechanism    string `json:"sasl.mechanism"`
	SaslUsername     string `json:"sasl.username"`
	SaslPassword     string `json:"sasl.password"`
}

// config should be a pointer to structure, if not, panic
func loadJsonConfig() *KafkaConfig {

	var config = &KafkaConfig{
		Topic:            os.Getenv("OrderTopic"),
		GroupId:          os.Getenv("OrderGroupId"),
		BootstrapServers: os.Getenv("BootstrapServers"),
		SecurityProtocol: os.Getenv("SecurityProtocol"),
		SaslMechanism:    os.Getenv("SaslMechanism"),
		SaslUsername:     os.Getenv("SaslUsername"),
		SaslPassword:     os.Getenv("SaslPassword"),
	}

	return config
}

func DoInitProducer(cfg *KafkaConfig) *kafka.Producer {
	fmt.Print("init kafka producer, it may take a few seconds to init the connection\n")
	//common arguments
	var kafkaconf = &kafka.ConfigMap{
		"api.version.request": "true",
		"message.max.bytes":   1000000,
		"linger.ms":           10,
		"retries":             30,
		"retry.backoff.ms":    1000,
		"acks":                "1"}
	kafkaconf.SetKey("bootstrap.servers", cfg.BootstrapServers)

	//公网必须使用证书 vpc 就不用了
	switch cfg.SecurityProtocol {
	case "plaintext":
		kafkaconf.SetKey("security.protocol", "plaintext")
	case "sasl_ssl":
		kafkaconf.SetKey("security.protocol", "sasl_ssl")
				//根据证书正式路径替换 /Users/victor/rfcwork/work/opentelemetry_demo/src/demo/order/config/ca-cert.pem
        kafkaconf.SetKey("ssl.ca.location", "/app/config/ca-cert.pem")
		kafkaconf.SetKey("sasl.username", cfg.SaslUsername)
		kafkaconf.SetKey("sasl.password", cfg.SaslPassword)
		kafkaconf.SetKey("sasl.mechanism", cfg.SaslMechanism)
		kafkaconf.SetKey("enable.ssl.certificate.verification", "false")
		kafkaconf.SetKey("ssl.endpoint.identification.algorithm", "None")
	case "sasl_plaintext":
		kafkaconf.SetKey("sasl.mechanism", "PLAIN")
		kafkaconf.SetKey("security.protocol", "sasl_plaintext")
		kafkaconf.SetKey("sasl.username", cfg.SaslUsername)
		kafkaconf.SetKey("sasl.password", cfg.SaslPassword)
		kafkaconf.SetKey("sasl.mechanism", cfg.SaslMechanism)
	default:
		panic(kafka.NewError(kafka.ErrUnknownProtocol, "unknown protocol", true))
	}

	producer, err := kafka.NewProducer(kafkaconf)
	if err != nil {
		panic(err)
	}
	fmt.Print("init kafka producer success\n")
	return producer
}

func DoInitConsumer(cfg *KafkaConfig) *kafka.Consumer {
	fmt.Print("init kafka consumer, it may take a few seconds to init the connection\n")
	//common arguments
	var kafkaconf = &kafka.ConfigMap{
		"api.version.request":       "true",
		"auto.offset.reset":         "earliest",
		"heartbeat.interval.ms":     3000,
		"session.timeout.ms":        30000,
		"max.poll.interval.ms":      120000,
		"fetch.max.bytes":           1024000,
		"max.partition.fetch.bytes": 256000}
	kafkaconf.SetKey("bootstrap.servers", cfg.BootstrapServers)
	kafkaconf.SetKey("group.id", cfg.GroupId)

	switch cfg.SecurityProtocol {
	case "plaintext":
		kafkaconf.SetKey("security.protocol", "plaintext")
	case "sasl_ssl":
		kafkaconf.SetKey("security.protocol", "sasl_ssl")
		//根据证书正式路径替换 /Users/victor/rfcwork/work/opentelemetry_demo/src/demo/order/config/ca-cert.pem
		kafkaconf.SetKey("ssl.ca.location", "/app/config/ca-cert.pem")
		kafkaconf.SetKey("sasl.username", cfg.SaslUsername)
		kafkaconf.SetKey("sasl.password", cfg.SaslPassword)
		kafkaconf.SetKey("sasl.mechanism", cfg.SaslMechanism)
		kafkaconf.SetKey("ssl.endpoint.identification.algorithm", "None")
		kafkaconf.SetKey("enable.ssl.certificate.verification", "false")
	case "sasl_plaintext":
		kafkaconf.SetKey("security.protocol", "sasl_plaintext")
		kafkaconf.SetKey("sasl.username", cfg.SaslUsername)
		kafkaconf.SetKey("sasl.password", cfg.SaslPassword)
		kafkaconf.SetKey("sasl.mechanism", cfg.SaslMechanism)

	default:
		panic(kafka.NewError(kafka.ErrUnknownProtocol, "unknown protocol", true))
	}

	consumer, err := kafka.NewConsumer(kafkaconf)
	if err != nil {
		panic(err)
	}
	fmt.Print("init kafka consumer success\n")
	return consumer
}

func SendKafka(msg LogisticsMsg) {
	config := loadJsonConfig()
	producer := DoInitProducer(config)

	defer producer.Close()
	// Delivery report handler for produced messages
	go func() {
		for e := range producer.Events() {
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

	marshal, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("SendKafka[] failed: %v\n", err.Error())
	}

	// Produce messages to topic (asynchronously)
	topic := config.Topic
	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          marshal,
	}, nil)

	// Wait for message deliveries before shutting down
	producer.Flush(150 * 1000)
}
func ServiceCallPost(ctx context.Context, host string, path string, data interface{}, res *util.Result) error {
	client := &util.HttpClient{
		Host: host,
	}
	err := client.Post(ctx, path, data, res)
	if err != nil {
		return err
	}
	if res.Code != 200 {
		s := fmt.Sprintf("%d %s %v", res.Code, res.Message, res.Data)
		return errors.New(s)
	}
	return nil
}

func ServiceCallGet(ctx context.Context, host string, path string, data map[string]interface{}, res *util.Result) error {

	client := &util.HttpClient{
		Host: host,
	}
	err := client.Get(ctx, path, data, res)
	if err != nil {
		return err
	}
	if res.Code != 200 {
		s := fmt.Sprintf("%d %s %v", res.Code, res.Message, res.Data)
		return errors.New(s)
	}
	return nil
}
