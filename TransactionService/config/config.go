package config

type KafkaConfig struct {
	Broker1Address string
	Topic          string
	ConsumerGroup  string
	// NumberOfConsumers int
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
			Broker1Address: "localhost:9092",
			Topic:          "CreationTransaction7",
			ConsumerGroup:  "TransactionService3",
			// NumberOfConsumers: 3,
		},
		DB: DBConfig{
			DbHost:     "localhost",
			DbPort:     5555,
			DbUser:     "docker",
			DbPassword: "docker",
			DbName:     "docker",
			SSLmode:    "disable",
		},
	}
	return config
}
