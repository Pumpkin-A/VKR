package config

type ServerConfig struct {
	Port int
}

type Config struct {
	Server ServerConfig
}

func New() Config {
	config := Config{
		Server: ServerConfig{
			Port: 8080,
		},
	}
	return config
}
