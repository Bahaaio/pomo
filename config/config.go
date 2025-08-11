// Package config loads, stores, and
// provides default values for work and break tasks.
package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type task struct {
	Duration time.Duration
	Then     []string
}

type config struct {
	Work      task
	Break     task
	AltScreen bool
}

var C config

func init() {
	viper.SetConfigName("pomo")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/pomo/")
	viper.AddConfigPath(".")

	viper.SetDefault("altScreen", true)

	viper.SetDefault("work", map[string]any{
		"duration": 25 * time.Minute,
	})

	viper.SetDefault("break", map[string]any{
		"duration": 5 * time.Minute,
	})
}

func LoadConfig() error {
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	log.Println("read config")

	err = viper.Unmarshal(&C)
	if err != nil {
		return err
	}
	log.Println("Unmarshaled into config struct")

	return nil
}

// Save writes the current configuration to file.
func Save() error {
	return viper.WriteConfig()
}
