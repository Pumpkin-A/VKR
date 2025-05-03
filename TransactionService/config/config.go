package config

type KafkaConfig struct {
	BootstrapServers string
	Topic            string
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
			BootstrapServers: "localhost:9092",
			Topic:            "CreationTransaction2",
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
