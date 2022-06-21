package config

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func InitLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func requireEnvVar(key string) string {
	x := os.Getenv(key)
	if x == "" {
		err := fmt.Errorf("surchable/internal/config: missing required environment variable %#v", key)
		err = errors.WithStack(err)
		log.Fatal().Stack().Err(err).Send()
	}
	return x
}

func envVarDefault(key string, defaultValue string) string {
	x := os.Getenv(key)
	if x == "" {
		return defaultValue
	}
	return x
}

/*[[[cog
import cog
from generateMappings import *
cog.outl(
	generate_golang_configuration(
		parse_configuration(
			load_raw_configuration(),
		),
	),
)
]]]*/
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
var Coordinator = struct {
	ServePort string
	ServeHost string
}{
	ServePort: envVarDefault("SURCHABLE_COORDINATOR_SERVE_PORT", "7200"),
	ServeHost: envVarDefault("SURCHABLE_COORDINATOR_SERVE_HOST", "0.0.0.0"),
}

// [[[end]]]
