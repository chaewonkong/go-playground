package main

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config 구조체
type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
		Mode string `mapstructure:"mode"`
	} `mapstructure:"server"`
	Database struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
	} `mapstructure:"database"`
}

func main() {
	v := viper.New()
	configPath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok {
		configPath = "config.yml"
	}

	var c Config
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	err := v.Unmarshal(&c)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("port: %s, secret: %s\n", c.Server.Port, c.Database.Password)
}
