package env

import (
	"os"
	"strconv"
)

func GetConfigPath() string {
	return getString("CONFIG_PATH", "./config/config.yaml")
}

func GetAPIKey() string {
	return getString("API_KEY", "-")
}

func GetDBUrl() string {
	return getString("DATABASE_URL", "-")
}

func getString(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return val
}

func getInt(key string, defaultValue int) int {
	strVal, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	intVal, err := strconv.ParseInt(strVal, 10, 64)
	if err != nil {
		return defaultValue
	}

	return int(intVal)
}
