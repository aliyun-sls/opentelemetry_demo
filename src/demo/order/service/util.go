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

/*
// writeByConn 基于Conn发送消息
func writeByConn(cfg *KafkaConfig, msg LogisticsMsg) {

	var dialer *kafka.Dialer

	switch cfg.SecurityProtocol {
	case "plaintext":
	case "sasl_ssl":
		// 1. 加载 SSL CA 证书（PEM 格式）
		caCert, err := os.ReadFile("/Users/victor/rfcwork/work/opentelemetry_demo/src/demo/order/config/ca-cert.pem")
		if err != nil {
			panic(fmt.Sprintf("读取 CA 证书失败: %v", err))
		}
		// 2. 创建 TLS 配置
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig := &tls.Config{
			//NameToCertificate: m,
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		}
		// 3. 创建 SASL 认证器
		saslMechanism := plain.Mechanism{
			Username: cfg.SaslUsername, // 阿里云 Kafka SASL 用户名
			Password: cfg.SaslPassword, // 阿里云 Kafka SASL 密码
		}
		// 4. 创建带有 TLS 和 SASL 的 Dialer
		dialer = &kafka.Dialer{
			DualStack:     true,
			SASLMechanism: saslMechanism,
			TLS:           tlsConfig,
		}
	case "sasl_plaintext":
		// 3. 创建 SASL 认证器
		saslMechanism := plain.Mechanism{
			Username: cfg.SaslUsername, // 阿里云 Kafka SASL 用户名
			Password: cfg.SaslPassword, // 阿里云 Kafka SASL 密码
		}
		// 4. 创建带有 TLS 和 SASL 的 Dialer
		dialer = &kafka.Dialer{
			DualStack:     true,
			SASLMechanism: saslMechanism,
		}
	default:
		panic(errors.New("unknown protocol"))

	}

	partition := 0

	// 连接至Kafka集群的Leader节点
	conn, err := dialer.DialLeader(context.Background(), "tcp", cfg.BootstrapServers, cfg.Topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	// 设置发送消息的超时时间
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// 发送消息
	marshal, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("SendKafka[] json.Marshal failed: %v\n", err.Error())
	}
	_, err = conn.WriteMessages(
		kafka.Message{Value: marshal},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
	time.Sleep(10000)
	// 关闭连接
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}*/
/*
// readByConn 连接至kafka后接收消息
func readByConn(cfg *KafkaConfig) {
	switch cfg.SecurityProtocol {
	case "plaintext":
	case "sasl_ssl":
	case "sasl_plaintext":
	default:
		panic(errors.New("unknown protocol"))

	}
	// 1. 加载 SSL CA 证书（PEM 格式）
	caCert, err := os.ReadFile("/Users/victor/rfcwork/work/opentelemetry_demo/src/demo/order/config/ca-cert.pem")
	if err != nil {
		panic(fmt.Sprintf("读取 CA 证书失败: %v", err))
	}
	// 2. 创建 TLS 配置
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	clientCert, err := tls.LoadX509KeyPair("client-cert.pem", "client-key.pem")
	if err != nil {
		log.Fatalf("Failed to load client certificate and key: %v", err)
	}
	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{clientCert}, // 设置客户端证书

	}
	// 3. 创建 SASL 认证器
	saslMechanism := plain.Mechanism{
		Username: cfg.SaslUsername, // 阿里云 Kafka SASL 用户名
		Password: cfg.SaslPassword, // 阿里云 Kafka SASL 密码
	}

	partition := 0
	// 4. 创建带有 TLS 和 SASL 的 Dialer
	dialer := &kafka.Dialer{
		DualStack:     true,
		SASLMechanism: saslMechanism,
		TLS:           tlsConfig,
	}

	// 连接至Kafka的leader节点
	conn, err := dialer.DialLeader(context.Background(), cfg.Topic, cfg.BootstrapServers, cfg.Topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	// 设置读取超时时间
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	// 读取一批消息，得到的batch是一系列消息的迭代器
	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	// 遍历读取消息
	b := make([]byte, 10e3) // 10KB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:n]))
	}

	// 关闭batch
	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}

	// 关闭连接
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}*/

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

func SendKafka(msg LogisticsMsg) error {
	config := loadJsonConfig()
	//writeByConn(config, msg)
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
		fmt.Printf("SendKafka[] json.Marshal failed: %v\n", err.Error())
		return err
	}

	// Produce messages to topic (asynchronously)
	topic := config.Topic
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          marshal,
	}, nil)

	if err != nil {
		fmt.Printf("SendKafka[] failed: %v\n", err.Error())
		return err
	}

	// Wait for message deliveries before shutting down
	producer.Flush(150 * 1000)
	return nil
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
