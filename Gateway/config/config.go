package config

type ServerConfig struct {
	Port int
}

type ClientConfig struct {
	ServerAddress string
}

type KafkaConfig struct {
	BootstrapServers string
	Topic            string
}

type Config struct {
	Server ServerConfig
	Kafka  KafkaConfig
	Client ClientConfig
}

func New() Config {
	config := Config{
		Server: ServerConfig{
			Port: 8080,
		},
		Kafka: KafkaConfig{
			BootstrapServers: "localhost:9092",
			Topic:            "ExternalTransactionOperations",
		},
		Client: ClientConfig{
			ServerAddress: "localhost:50051",
		},
	}
	return config
}
