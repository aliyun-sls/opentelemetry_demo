package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
	"os/signal"
	"sls-mall-go/common/util"
	"syscall"
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
		kafkaconf.SetKey("ssl.ca.location", "/usr/src/app/config/ca-cert.pem")
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
		"auto.offset.reset":         "latest",
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
		kafkaconf.SetKey("ssl.ca.location", "/usr/src/app/config/ca-cert.pem")
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

func Consumerkafka() {
	cfg := loadJsonConfig()
	consumer := DoInitConsumer(cfg)

	consumer.SubscribeTopics([]string{cfg.Topic}, nil)

	// 处理系统信号以优雅关闭
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			select {
			case sig := <-sigchan:
				fmt.Printf("Caught signal %v: terminating\n", sig)
				consumer.Close()
				return
			default:
				fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
				var logisticsMsg = &LogisticsMsg{}
				err = json.Unmarshal(msg.Value, logisticsMsg)
				if err != nil {
					fmt.Printf("Consumer Unmarshal error: %v (%v)\n", err, msg)
				}
				switch logisticsMsg.Action {
				case Create:
					CreateMsg(logisticsMsg)
				}
			}
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	consumer.Close()
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
