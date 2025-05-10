package config

type KafkaConfig struct {
	Broker1Address                     string
	InternalPaymentResultTopic         string
	InternalTransactionOperationsTopic string
	ConsumerGroup                      string
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
			BaseURL: "http://localhost:9999",
		},
		Kafka: KafkaConfig{
			Broker1Address:                     "localhost:9092",
			InternalPaymentResultTopic:         "InternalPaymentResult",
			InternalTransactionOperationsTopic: "InternalTransactionOperations",
			ConsumerGroup:                      "Billing",
		},
	}
	return config
}
