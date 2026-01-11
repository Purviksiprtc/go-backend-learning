package config

import "os"

func GetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("Missing required env var: " + key)
	}
	return val
}
