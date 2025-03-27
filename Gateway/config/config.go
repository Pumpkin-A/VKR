package config

type ServerConfig struct {
	Port int
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
	Server ServerConfig
	DB     DBConfig
}

func New() Config {
	config := Config{
		Server: ServerConfig{
			Port: 8080,
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
