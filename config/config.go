// Package config loads, stores, and
// provides default values for work and break tasks.
package config

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Notification struct {
	Enabled bool
	Title   string
	Message string
	Icon    string
}

type Task struct {
	Title        string
	Duration     time.Duration
	Then         [][]string
	Notification Notification
}

type Config struct {
	Work          Task
	Break         Task
	AskToContinue bool
}

var (
	//go:embed pomo.png
	Icon []byte
	C    Config

	defaultConfig = map[string]any{
		"askToContinue": true,
		"work": Task{
			Duration: 25 * time.Minute,
			Title:    "work session",
			Notification: Notification{
				Enabled: true,
				Title:   "work finished 🎉",
				Message: "time to take a break",
			},
		},
		"break": Task{
			Duration: 5 * time.Minute,
			Title:    "break session",
			Notification: Notification{
				Enabled: true,
				Title:   "break over 😴",
				Message: "back to work!",
			},
		},
	}
)

func setup() {
	viper.SetConfigName("pomo")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")
	if configDir, err := os.UserConfigDir(); err == nil {
		viper.AddConfigPath(filepath.Join(configDir, "pomo"))
	} else {
		log.Println("could not get user config dir:", err)
	}

	for k, v := range defaultConfig {
		viper.SetDefault(k, v)
	}
}

func LoadConfig() error {
	setup()
	log.Println("loading config")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	log.Println("read config:", viper.ConfigFileUsed())

	err = viper.Unmarshal(&C)
	if err != nil {
		return err
	}

	log.Println("Unmarshaled config:", C)
	return nil
}
