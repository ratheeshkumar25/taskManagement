package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Database_url   string `mapstructure:"DATABASE_URL"`
	REDISHOST      string `mapstructure:"REDISHOST"`
	REDIS_PASSWORD string `mapstructure:"REDIS_PASSWORD"`
	SECERETKEY     string `mapstructure:"JWTKEY"`
	PORT           string `mapstructure:"PORT"`
}

func LoadConfig() *Config {
	var config Config

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("No .env file found. Using environment variables.")
	}

	viper.AutomaticEnv()

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	//log.Println("From Viper directly:", viper.GetString("DATABASE_URL"))

	return &config
}
