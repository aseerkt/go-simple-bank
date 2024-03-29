package utils

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBUser        string `mapstructure:"DB_USER"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_NAME"`
	DBPort        string `mapstructure:"DB_PORT"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	DBUrl         string `mapstructure:"DB_URL"`
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
