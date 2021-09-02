package config

import (
	http2 "broken-link-checker/app/internal/delivery/http"
	http_test2 "broken-link-checker/app/internal/delivery/http_test"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Server     http2.Config
	ServerTest http_test2.Config
}

var (
	config Config
	once   sync.Once
)

// Get reads config from environment and config.yml. Once.
func Get() *Config {
	once.Do(func() {
		configPath := "./configs"

		if err := checkConfigPath(configPath); err != nil {
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

func checkConfigPath(path string) error {
	var (
		pathConfig        = path + "/config.yml"
		pathExampleConfig = path + "/config_example.yml"
	)

	if _, err := os.Stat(pathConfig); err != nil {
		if os.IsNotExist(err) {
			// config doesn't exist
			if _, err := os.Stat(pathExampleConfig); err == nil {
				// config_example exists. Creating a config based on the example

				if copyErr := copyFileContents(pathExampleConfig, pathConfig); copyErr != nil {
					return copyErr
				}
			} else {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cErr := out.Close()
		if cErr == nil {
			err = cErr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()

	return
}
