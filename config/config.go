package config

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewConfig(configFilePath string) {
	viper.SetConfigName("gladius-guardian")
	viper.AddConfigPath(configFilePath)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("Fatal error reading config file: %s \n", err))
	}

	ConfigOption("networkdExecutable", "gladius-networkd")
	ConfigOption("controldExectuable", "gladius-controld")

}

func ConfigOption(key string, default_value interface{}) string {
	viper.SetDefault(key, default_value)

	return key
}
