package bootstrap

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	configApp "mapexIam/src/shared/configuration/application"
)

// InitConfig initializes the application configuration from environment and defaults.
func InitConfig() {
	config.InitConfig(configApp.DefaultConfiguration)
}

// InitLogger initializes the structured logger based on service configuration.
func InitLogger() {
	serviceName, _ := config.GetStringValue("service_name")
	serviceVersion, _ := config.GetStringValue("service_version")
	goEnv, _ := config.GetStringValue("go_env")

	// Set default log level based on environment
	level := logger.DebugLevel
	if goEnv != "development" && goEnv != "dev" {
		level = logger.InfoLevel
	}

	// LOG_LEVEL env overrides the GO_ENV-based default
	logLevel, _ := config.GetStringValue("log_level")
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
	logger.Info("[APP:BOOTSTRAP] MapexOS Service starting")
	logger.Info("[APP:BOOTSTRAP] Logger initialized")
}
