package app

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// env var에서 config를 읽어온다
func ReadConfig(cfg any) error {
	configPath := os.Getenv(envConfigFilePath) // config.yml path는 environ으로 관리

	v := viper.New()
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
