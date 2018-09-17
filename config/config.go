package config

import (
	"fmt"

	gconfig "github.com/gladiusio/gladius-utils/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func SetupConfig(configFilePath string) {
	viper.SetConfigName("gladius-guardian")
	viper.AddConfigPath(configFilePath)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("Fatal error reading config file: %s \n", err))
	}

	ConfigOption("NetworkdExecutable", "gladius-networkd")
	ConfigOption("ControldExectuable", "gladius-controld")

	// Setup gladius base for the various services
	base, err := gconfig.GetGladiusBase()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Warn("Couldn't get Gladius base")
	}
	ConfigOption("DefaultEnvironment", []string{"GLADIUSBASE=" + base})
}

func ConfigOption(key string, defaultValue interface{}) string {
	viper.SetDefault(key, defaultValue)

	return key
}
