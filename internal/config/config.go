package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// struct of type config
type Config struct {
	DB_HOST     string `mapstructure:"DB_HOST"`
	DB_PORT     string `mapstructure:"DB_PORT"`
	DB_USER     string `mapstructure:"DB_USER"`
	DB_PASSWORD string `mapstructure:"DB_PASSWORD"`
	DB_NAME     string `mapstructure:"DB_NAME"`
	DB_SSLMODE  string `mapstructure:"DB_SSLMODE"`
	APP_ENV     string `mapstructure:"APP_ENV"`
	APP_PORT    string `mapstructure:"APP_PORT"`
	JWT_KEY     string `mapstructure:"JWT_KEY"`
}

// Function used to load config accepts no parameters and returns a pointer to Config struct and error(if any)
func LoadConfig() (*Config, error) {
	// create empty struct of type Config
	config := &Config{}
	env := "local"
	envConfigFileName := fmt.Sprintf(".env.%s", env)
	viper.AutomaticEnv()
	viper.AddConfigPath("./.secrets")
	viper.SetConfigName(envConfigFileName)
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Config file not found. Using env variables")
		} else {
			return nil, fmt.Errorf("failed to read config file:  %w", err)
		}
	}

	err = viper.Unmarshal(&config)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall config :%w", err)
	}
	fmt.Println("config: ", config)

	return config, nil

}
