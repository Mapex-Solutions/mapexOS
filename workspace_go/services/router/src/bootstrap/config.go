package bootstrap

import (
	"fmt"

	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	configApp "router/src/shared/configuration/application"
)

// InitConfig initializes the application configuration from environment and defaults.
func InitConfig() {
	config.InitConfig(configApp.DefaultConfiguration)
}

// InitLogger initializes the structured logger based on service configuration.
//
// Log level priority: LOG_LEVEL env var > GO_ENV-based default.
// Valid LOG_LEVEL values: debug, info, warn, error, silent.
func InitLogger() {
	serviceName, _ := config.GetStringValue("service_name")
	serviceVersion, _ := config.GetStringValue("service_version")
	goEnv, _ := config.GetStringValue("go_env")
	logLevel, _ := config.GetStringValue("log_level")

	// Resolve log level: explicit LOG_LEVEL takes priority
	level := logger.InfoLevel
	if goEnv == "development" || goEnv == "dev" {
		level = logger.DebugLevel
	}

	switch logLevel {
	case "debug":
		level = logger.DebugLevel
	case "info":
		level = logger.InfoLevel
	case "warn":
		level = logger.WarnLevel
	case "error":
		level = logger.ErrorLevel
	case "silent":
		level = logger.DisabledLevel
	}

	logger.InitLogger(logger.LoggerOptions{
		ServiceName:    serviceName,
		ServiceVersion: serviceVersion,
		Environment:    goEnv,
		Level:          level,
	})
	logger.Info("[APP:BOOTSTRAP] Router Service starting")
	logger.Info(fmt.Sprintf("[APP:BOOTSTRAP] Logger initialized: goEnv=%s, logLevel=%s, resolvedLevel=%d", goEnv, logLevel, level))
}
