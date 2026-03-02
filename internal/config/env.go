package config

import (
	"log"
	"os"
	"strconv"
)

func getEnv(key, defaultValue string) string {
	log.Println(key)
	value := os.Getenv(key)
	log.Println("Test", value)
	if value == "" {
		return defaultValue
	}
	return value

}

func getEnvInt(key string, defaultValue int) int {

	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)

	if err != nil {
		return defaultValue
	}

	return value

}
