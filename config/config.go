package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config contains env variables
type Config struct {
	ENVIRONMENT string `mapstructure:"environment"`
	DEBUG       bool   `mapstructure:"debug"`
	SERVER      server
	DB          db
}

type server struct {
	HOST string `mapstructure:"host"`
	PORT int    `mapstructure:"port"`
}

type db struct {
	HOST     string `mapstructure:"host"`
	PORT     int    `mapstructure:"port"`
	USERNAME string `mapstructure:"username"`
	PASSWORD string `mapstructure:"password"`
	DBNAME   string `mapstructure:"db_name"`
}

// Conf is global configuration file
var Conf Config

// InitConfig initialze the Conf
func InitConfig() {
	config := initViper()
	Conf = *config
}

func initViper() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Panicf("Unable to decode into struct, %v", err)
	}
	return &config
}
