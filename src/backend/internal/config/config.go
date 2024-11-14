package config

import "github.com/joho/godotenv"

type Config struct {
	DBConfig        DBConfig
	WoWClientID     string
	WoWClientSecret string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{}, nil
}