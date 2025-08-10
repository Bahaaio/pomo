// Package config loads, stores, and
// provides default values for work and break tasks.
package config

import (
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
		"then": []string{
			"notify-send 'Work Finished!' 'Time to take a break ☕'",
		},
	})

	viper.SetDefault("break", map[string]any{
		"duration": 5 * time.Minute,
		"then": []string{
			"notify-send 'Break Over'",
		},
	})

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("config file not found")
		} else {
			panic("could not read config file")
		}
	}

	err = viper.Unmarshal(&C)
	if err != nil {
		panic("invalid config file")
	}
}

// Save writes the current configuration to file.
func Save() error {
	return viper.WriteConfig()
}
