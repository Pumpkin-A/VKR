package config

type KafkaConfig struct {
	Broker1Address string
	Topic          string
	ConsumerGroup  string
}

type BankExampleClientConfig struct {
	BaseURL string
}

type Config struct {
	Client BankExampleClientConfig
	Kafka  KafkaConfig
}

func New() Config {
	config := Config{
		Client: BankExampleClientConfig{
			BaseURL: "http://localhost:9876",
		},
		Kafka: KafkaConfig{
			Broker1Address: "localhost:9092",
			Topic:          "TransactionFinishStatus",
			ConsumerGroup:  "Statistics",
		},
	}
	return config
}
