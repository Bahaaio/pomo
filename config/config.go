// Package config loads, stores, and
// provides default values for work and break tasks.
package config

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	log.Println("setting default config values")
	for k, v := range defaultConfig {
		viper.SetDefault(k, v)
	}
}

func LoadConfig() error {
	setup()
	log.Println("loading config")

	// fall back to defaults if no config file is found
	if err := viper.ReadInConfig(); err != nil {
		log.Println("no config file found, using defaults:", err)
	} else {
		log.Println("read config:", viper.ConfigFileUsed())
	}

	err := viper.Unmarshal(&C)
	if err != nil {
		return err
	}
	log.Println("Unmarshaled config:", C)

	if C.Work.Notification.Icon, err = expandPath(C.Work.Notification.Icon); err != nil {
		log.Println("failed to expand Work Notification icon path:", err)
	}

	if C.Break.Notification.Icon, err = expandPath(C.Break.Notification.Icon); err != nil {
		log.Println("failed to expand Break Notification icon path:", err)
	}

	return nil
}

// expands tilde to the user's home directory
func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %v", err)
		}

		expandedPath := filepath.Join(homeDir, path[2:])
		log.Printf("expanding path: %s to %s\n", path, expandedPath)

		return expandedPath, nil
	}

	return path, nil
}
