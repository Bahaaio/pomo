// Package config loads, stores, and
// provides default values for work and break tasks.
package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Task struct {
	Duration time.Duration
	Then     []string
}

type Config struct {
	Work       Task
	Break      Task
	FullScreen bool
}

var C Config

func init() {
	viper.SetConfigName("pomo")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/pomo/")
	viper.AddConfigPath(".")

	viper.SetDefault("fullScreen", true)

	viper.SetDefault("work", map[string]any{
		"duration": 25 * time.Minute,
	})

	viper.SetDefault("break", map[string]any{
		"duration": 5 * time.Minute,
	})
}

func LoadConfig() error {
	log.Println("loading config")

	viper.ReadInConfig()
	log.Println("read config")

	err := viper.Unmarshal(&C)
	if err != nil {
		return err
	}
	log.Println("Unmarshaled config")
	log.Println("config:", C)

	return nil
}
