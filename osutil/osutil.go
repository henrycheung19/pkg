package osutil

import (
	"encoding/json"
	"os"
	"strconv"
)

// GetEnvStr return a string if the key is present in the environment. If the key is not present, then return the inputted default value.
func GetEnvStr(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return def
}

// GetEnvInt return an int if the key is present in the environment. If the key is not present, then return the inputted default value.
func GetEnvInt(key string, def int) int {
	if value := os.Getenv(key); value != "" {
		i, err := strconv.Atoi(value)
		if err != nil {
			return def
		}
		return i
	}
	return def
}

// GetEnvBool return a boolean if the key is present in the environment. If the key is not present, then return the inputted default value.
func GetEnvBool(key string, def bool) bool {
	if value := os.Getenv(key); value == "true" {
		return true
	}
	return def
}

// GetEnvStrArray return an array if the key is present in the environment. If the key is not present, then return the inputted default value.
func GetEnvStrArray(key string, def []string) []string {
	if value := os.Getenv(key); value != "" {
		var arr []string
		json.Unmarshal([]byte(key), &arr)
		return arr
	}
	return def
}
