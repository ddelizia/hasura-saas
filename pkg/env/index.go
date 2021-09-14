package env

import (
	"fmt"
	"net/url"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var firstLoad = false

func loadEnvs() {
	logrus.Info("loading configuration")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../") // for testing porpouse
	viper.AddConfigPath("../")    // for testing porpouse
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func lazyLoadEnv() {
	if !firstLoad {
		firstLoad = true
		loadEnvs()
	}
}

// Get configuration variable as string
func GetString(key string) string {
	lazyLoadEnv()
	return viper.GetString(key)
}

// Get configuration variable as url
func GetUrl(key string) *url.URL {
	lazyLoadEnv()
	redirectUrl := GetString(key)
	url, err := url.Parse(redirectUrl)
	if err != nil {
		logrus.Error(key, " environment variable is not an url: ", redirectUrl)
		os.Exit(1)
	}
	return url
}
