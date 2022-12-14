package config

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Val Val
var Val Config

// Config Config
type Config struct {
	Mode string `mapstructure:"MODE"`
	//Port string `mapstructure:"PORT"`

	GoogleSecretKey string `mapstructure:"GOOGLE_SECRET_KEY"`
	GoogleClientID  string `mapstructure:"GOOLE_CLIENT_ID"`
}

// Init Init
func Init() {
	// 讀config.yaml
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %v ", err))
	}
	if err := viper.Unmarshal(&Val); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

	log.WithFields(log.Fields{
		"val": Val,
	}).Info("config loaded")
}
