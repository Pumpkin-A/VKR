package config

type ServerConfig struct {
	Port int
}

type KafkaConfig struct {
	BootstrapServers string
	Topic            string
}

type Config struct {
	Server ServerConfig
	Kafka  KafkaConfig
}

func New() Config {
	config := Config{
		Server: ServerConfig{
			Port: 8080,
		},
		Kafka: KafkaConfig{
			BootstrapServers: "localhost:9092",
			Topic:            "CreationTransaction2",
		},
	}
	return config
}
