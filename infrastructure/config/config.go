package config

import "os"

type DBConfig struct {
	ConnString string
}

type Config struct {
	DBConfig DBConfig
	Secret   string
	Port     string
}

func ReadEnv() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}
	return Config{
		Port: port,
		DBConfig: DBConfig{
			ConnString: os.Getenv("DB_CONN"),
		},
		Secret: os.Getenv("SECRET"),
	}
}
