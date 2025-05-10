package config

import (
	"os"
	"strings"
)

func DetectEnv() string {
	env := os.Getenv("APP_ENV")
	switch strings.ToLower(env) {
	case "prod", "production":
		return "prod"
	case "dev", "development":
		return "dev"
	default:
		return ""
	}
}
