package common

import (
	"log"
	"os"
	"strconv"
	"time"
)

func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid value for %s: %v. Using default: %d", key, err, defaultValue)
		return defaultValue
	}

	return time.Duration(value)
}
