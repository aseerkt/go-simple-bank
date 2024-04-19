package utils

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Mode          string `mapstructure:"GIN_MODE"`
	DBDriver      string `mapstructure:"DB_DRIVER"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	DBUrl         string `mapstructure:"DB_URL"`
	MigrateUrl    string `mapstructure:"MIGRATE_URL"`
}

func LoadConfig(path string) Config {
	var config Config

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("unable to read config: ", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("unable to unmarshal config: ", err)
	}

	return config
}
