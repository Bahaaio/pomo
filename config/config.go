// Package config loads, stores, and
// provides default values for work and break tasks.
package config

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Bahaaio/pomo/ui/ascii"
	"github.com/Bahaaio/pomo/ui/colors"
	"github.com/spf13/viper"
)

const (
	AppName    = "pomo"
	ConfigFile = "pomo.yml"
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

	DefaultConfig = map[string]any{
		"askToContinue": true,
		"asciiArt": map[string]any{
			"enabled": true,
			"font":    ascii.DefaultFont,
			"color":   colors.TimerFg,
		},
		"work": map[string]any{
			"duration": 25 * time.Minute,
			"title":    "work session",
			"notification": map[string]any{
				"enabled": true,
				"title":   "work finished 🎉",
				"message": "time to take a break!",
			},
		},
		"break": map[string]any{
			"duration": 5 * time.Minute,
			"title":    "break session",
			"notification": map[string]any{
				"enabled": true,
				"title":   "break over 😴",
				"message": "back to work!",
			},
		},
	}
)

func Setup() {
	viper.SetConfigName(AppName)
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")

	if configDir, err := getConfigDir(); err == nil {
		viper.AddConfigPath(configDir)
	} else {
		log.Println("could not get user config dir:", err)
	}

	log.Println("setting default config values")
	setDefaults()
}

func LoadConfig() error {
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

func setDefaults() {
	for key, value := range DefaultConfig {
		viper.SetDefault(key, value)
	}
}

// expands tilde to the user's home directory
func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %v", err)
		}

		return filepath.Join(homeDir, path[2:]), nil
	}

	return path, nil
}

// returns the config directory for the app
func getConfigDir() (string, error) {
	var dir string

	// on linux and macOS, use ~/.config
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		dir = os.Getenv("HOME")
		if dir == "" {
			return "", errors.New("$HOME is not defined")
		}

		dir = filepath.Join(dir, ".config")
	} else {
		// on other OSes, use the standard user config directory
		var err error
		dir, err = os.UserConfigDir()
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(dir, AppName), nil
}
