package config

import (
	"fmt"
	"os"
)

func requireEnvVar(key string) string {
	x := os.Getenv(key)
	if x == "" {
		panic(fmt.Errorf("surchable/internal/config: missing required environment variable %#v", key))
	}
	return x
}

var DB = struct {
	DatabaseName string
	User         string
	Password     string
	Host         string
}{
	DatabaseName: requireEnvVar("SURCHABLE_DB_DATABASE_NAME"),
	User:         requireEnvVar("SURCHABLE_DB_USER"),
	Password:     requireEnvVar("SURCHABLE_DB_PASSWORD"),
	Host:         requireEnvVar("SURCHABLE_DB_HOST"),
}
