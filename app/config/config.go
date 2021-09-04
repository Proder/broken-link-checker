package config

import (
	"broken-link-checker/app/internal/delivery/http"
	"broken-link-checker/app/internal/delivery/http_test"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Server     http.Config
	ServerTest http_test.Config
}

var (
	config Config
	once   sync.Once
)

// Get reads config from environment and config.yml. Once.
func Get() *Config {
	once.Do(func() {
		configPath := "./configs"

		if _, err := os.Stat(configPath + "/config.yml"); err != nil {
			log.Fatal(err)
			return
		}

		viper.AddConfigPath(configPath)
		viper.SetConfigName("config")

		viper.AllowEmptyEnv(true)
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()

		// Read config.yml
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal(err)
			return
		}

		if err := viper.Unmarshal(&config); err != nil {
			log.Fatal(err)
			return
		}
	})

	return &config
}
