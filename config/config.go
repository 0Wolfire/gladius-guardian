package config

import (
	"fmt"
	"strings"

	"github.com/gladiusio/gladius-common/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func SetupConfig(configFilePath string) {
	viper.SetConfigName("gladius-guardian")
	viper.AddConfigPath(configFilePath)

	// Setup env variable handling
	viper.SetEnvPrefix("GUARDIAN")
	r := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(r)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Warn(fmt.Errorf("error reading config file: %s, using defaults", err))
	}

	ConfigOption("NetworkdExecutable", "gladius-networkd")
	ConfigOption("ControldExecutable", "gladius-controld")

	// Setup gladius base for the various services
	base, err := utils.GetGladiusBase()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Warn("Couldn't get Gladius base")
	}
	// Add a default environment so that we can set the gladius base of our sub
	// processes
	ConfigOption("DefaultEnvironment", []string{"GLADIUSBASE=" + base})

	ConfigOption("MaxLogLines", 1000) // Max number of log lines to keep in ram for each service

	// Setup logging level
	switch loglevel := viper.GetString("LogLevel"); loglevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func ConfigOption(key string, defaultValue interface{}) string {
	viper.SetDefault(key, defaultValue)

	return key
}
