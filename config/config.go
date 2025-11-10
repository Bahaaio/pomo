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

	"github.com/Bahaaio/pomo/ui/ascii"
	"github.com/Bahaaio/pomo/ui/colors"
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

type ASCIIArt struct {
	Enabled bool
	Font    string
	Color   string
}

type Config struct {
	Work          Task
	Break         Task
	AskToContinue bool
	ASCIIArt      ASCIIArt
}

var (
	//go:embed pomo.png
	Icon []byte
	C    Config
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

	viper.SetDefault("askToContinue", true)

	viper.SetDefault("asciiArt", map[string]any{
		"enabled": true,
		"font":    ascii.DefaultFont,
		"color":   colors.TimerFg,
	})

	viper.SetDefault("work", map[string]any{
		"duration": 25 * time.Minute,
		"title":    "work session",
		"notification": map[string]any{
			"enabled": true,
			"title":   "work finished 🎉",
			"message": "time to take a break",
		},
	})

	viper.SetDefault("break", map[string]any{
		"duration": 5 * time.Minute,
		"title":    "break session",
		"notification": map[string]any{
			"enabled": true,
			"title":   "break over 😴",
			"message": "back to work!",
		},
	})
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
