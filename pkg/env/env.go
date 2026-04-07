package env

import "os"

func GetConfigPath() string {
	return getString("CONFIG_PATH", "./config/config.yaml")
}

func getString(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return val
}
