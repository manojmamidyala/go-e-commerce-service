package config

import (
	"log"

	"github.com/spf13/viper"
)

type envConfigs struct {
	Environment   string `mapstructure:"ENVIRONMENT"`
	HttpPort      int    `mapstructure:"HTTP_PORT"`
	AuthSecret    string `mapstructure:"AUTH_SECRET"`
	DatabaseURL   string `mapstructure:"DATABASE_URL"`
	RedisURL      string `mapstructure:"REDIS_URI"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       string `mapstructure:"REDIS_DB"`
}

var appEnvVariables envConfigs

// Call to load the variables from env
func LoadEnvVariables() *envConfigs {
	// Tell viper the path/location of your env file. If it is root just add "."
	viper.AddConfigPath("./config/")

	// Tell viper the name of your file
	viper.SetConfigName("app")

	// Tell viper the type of your file
	viper.SetConfigType("env")

	// Viper reads all the variables from env file and log error if any found
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	// Viper unmarshals the loaded env varialbes into the struct
	if err := viper.Unmarshal(&appEnvVariables); err != nil {
		log.Fatal(err)
	}

	return &appEnvVariables
}

func GetConfig() *envConfigs {
	return &appEnvVariables
}
