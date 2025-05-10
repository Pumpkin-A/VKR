package config

type KafkaConfig struct {
	Broker1Address                     string
	ExternalTransactionOperationsTopic string
	InternalTransactionOperationsTopic string
	ConsumerGroup                      string
}

type DBConfig struct {
	DbHost     string
	DbPort     int
	DbUser     string
	DbPassword string
	DbName     string
	SSLmode    string
}

type Config struct {
	// Server ServerConfig
	Kafka KafkaConfig
	DB    DBConfig
}

func New() Config {
	config := Config{
		// Server: ServerConfig{
		// 	Port: 8080,
		// },
		Kafka: KafkaConfig{
			Broker1Address:                     "localhost:9092",
			ExternalTransactionOperationsTopic: "ExternalTransactionOperations",
			InternalTransactionOperationsTopic: "InternalTransactionOperations",
			ConsumerGroup:                      "TransactionService",
		},
		DB: DBConfig{
			DbHost:     "localhost",
			DbPort:     5432,
			DbUser:     "postgres",
			DbPassword: "postgres",
			DbName:     "vkr",
			SSLmode:    "disable",
		},
	}
	return config
}
