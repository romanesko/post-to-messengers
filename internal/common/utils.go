package common

import (
	"log"
	"os"
	"strings"
)

func GetEnv(name string, defaultValue ...string) string {
	value, exists := os.LookupEnv(name)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		log.Fatalf("Environment variable '%s' is not set and no default value provided", name)
	}
	return value
}

func GetEnvBool(envVar string) bool {
	value, exists := os.LookupEnv(envVar)
	if !exists {
		return false
	}
	normalizedValue := strings.TrimSpace(strings.ToLower(value))
	return normalizedValue != "" && normalizedValue != "0" && normalizedValue != "false"
}
