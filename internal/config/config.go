package config

import (
	"broken-link-checker/internal/delivery/http"
	"broken-link-checker/internal/delivery/http_test"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server     http.Config
	ServerTest http_test.Config
}

// Get reads config from environment and config.yml. Once.
func Get() (Config, error) {
	cnf := Config{}
	configPath := "./configs"

	if _, err := os.Stat(configPath + "/config.yml"); err != nil {
		return cnf, fmt.Errorf("os.Stat failed: %w", err)
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")

	viper.AllowEmptyEnv(true)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read config.yml
	if err := viper.ReadInConfig(); err != nil {
		return cnf, fmt.Errorf("viper.ReadInConfig failed: %w", err)
	}

	if err := viper.Unmarshal(&cnf); err != nil {
		return cnf, fmt.Errorf("viper.Unmarshal failed: %w", err)
	}

	return cnf, nil
}
