package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	DatabaseConnectionString string
	Port                     int
}

func FromEnv() (*Config, error) {
	databaseConnectionString := os.Getenv("DATABASE_URL")
	if databaseConnectionString == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	port := 4000
	portStr := os.Getenv("PORT")
	if portStr != "" {
		portInt, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, errors.New("PORT must be an integer")
		}
		port = portInt
	}

	return &Config{
		DatabaseConnectionString: databaseConnectionString,
		Port:                     port,
	}, nil
}
