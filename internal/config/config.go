package config

type Config struct {
	ServerPort string
}

func LoadConfig() *Config {

	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}
