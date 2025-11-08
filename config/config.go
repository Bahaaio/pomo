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
)

type TaskType int

const (
	WorkTask TaskType = iota
	BreakTask
)

func (t TaskType) GetTask() *Task {
	if t == BreakTask {
		return &C.Break
	}
	return &C.Work
}

func (t TaskType) Opposite() TaskType {
	if t == WorkTask {
		return BreakTask
	}
	return WorkTask
}

func init() {
	viper.SetConfigName("pomo")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")
	if configDir, err := os.UserConfigDir(); err == nil {
		viper.AddConfigPath(filepath.Join(configDir, "pomo"))
	}

	viper.SetDefault("askToContinue", true)

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
	log.Println("loading config")

	_ = viper.ReadInConfig()
	log.Println("read config")

	err := viper.Unmarshal(&C)
	if err != nil {
		return err
	}
	log.Println("Unmarshaled config")
	log.Println("config:", C)

	return nil
}
