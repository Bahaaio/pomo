package config

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

func TestLoadConfig(t *testing.T) {
	testConfig := `
askToContinue: true
work:
  duration: 31m
  title: custom work
`

	setupViper()
	writeAndLoadConfig(t, testConfig)

	assert.True(t, C.AskToContinue, "AskToContinue should be true")
	assert.Equal(t, 31*time.Minute, C.Work.Duration, "Work duration should be 31 minutes")
	assert.Equal(t, "custom work", C.Work.Title, "Work title should be 'custom work'")
}

func TestLoadConfigDefaults(t *testing.T) {
	setupViper()

	// load config with no config file present (should use defaults)
	writeAndLoadConfig(t, "")

	// compare loaded config with expected defaults
	defaults := getDefaultConfig()
	assertConfigMatches(t, defaults, C)
}

func TestLoadConfigPartialUpdate(t *testing.T) {
	partialConfig := `
askToContinue: true
work:
  duration: 30m
  title: custom work
asciiArt:
  color: "#FF0000"
`

	setupViper()
	writeAndLoadConfig(t, partialConfig)

	// test overridden values
	assert.True(t, C.AskToContinue, "AskToContinue should be true")
	assert.Equal(t, 30*time.Minute, C.Work.Duration, "Work duration should be 30 minutes")
	assert.Equal(t, "custom work", C.Work.Title, "Work title should be 'custom work'")
	assert.Equal(t, "#FF0000", C.ASCIIArt.Color, "ASCII art color should be '#FF0000'")

	defaults := getDefaultConfig()

	// test default values
	assert.Equal(t, defaults.Break.Duration, C.Break.Duration)
	assert.Equal(t, defaults.Break.Title, C.Break.Title)
	assert.Equal(t, defaults.Work.Notification.Title, C.Work.Notification.Title)
	assert.Equal(t, defaults.ASCIIArt.Enabled, C.ASCIIArt.Enabled)
	assert.Equal(t, defaults.ASCIIArt.Font, C.ASCIIArt.Font)
}

func TestLoadConfigThenCommands(t *testing.T) {
	configYAML := `
work:
  then:
    - ["echo", "Work session completed"]
    - ["osascript", "-e", "display notification \"Break time!\""]
    - ["python", "~/scripts/work-done.py"]
`

	setupViper()
	writeAndLoadConfig(t, configYAML)

	// Test work then commands
	expectedThen := [][]string{
		{"echo", "Work session completed"},
		{"osascript", "-e", "display notification \"Break time!\""},
		{"python", "~/scripts/work-done.py"},
	}
	assert.Equal(t, expectedThen, C.Work.Then, "Work then commands should match")
}

func TestLoadConfigAllFieldsComprehensive(t *testing.T) {
	configYAML := `
askToContinue: true
asciiArt:
  enabled: true
  font: ansi
  color: "#FF5733"
work:
  duration: 45m
  title: Deep work session
  notification:
    enabled: true
    title: Work Complete! 🎉
    message: Time for a well-deserved break
    icon: C:\User\image.svg
  then:
    - [echo, "work completed"]
    - [notify-send, "Break time!"]
break:
  duration: 15m
  title: Relaxation break
  notification:
    enabled: false
    title: Break Over! 😴
    message: Back to productive work
    icon: /abs/path/break-icon.png
  then:
    - ["echo", "break finished"]
`

	setupViper()
	writeAndLoadConfig(t, configYAML)

	// main config
	assert.True(t, C.AskToContinue, "AskToContinue should be true")

	// ASCII art
	assert.True(t, C.ASCIIArt.Enabled, "ASCII art should be enabled")
	assert.Equal(t, "ansi", C.ASCIIArt.Font, "ASCII art font should be 'ansi'")
	assert.Equal(t, "#FF5733", C.ASCIIArt.Color, "ASCII art color should be '#FF5733")

	// work task
	assert.Equal(t, 45*time.Minute, C.Work.Duration, "Work duration should be 45 minutes")
	assert.Equal(t, "Deep work session", C.Work.Title, "Work title should match")
	expectedWorkThen := [][]string{{"echo", "work completed"}, {"notify-send", "Break time!"}}
	assert.Equal(t, expectedWorkThen, C.Work.Then, "Work then commands should match")

	// work notification
	assert.True(t, C.Work.Notification.Enabled, "Work notification should be enabled")
	assert.Equal(t, "Work Complete! 🎉", C.Work.Notification.Title, "Work notification title should match")
	assert.Equal(t, "Time for a well-deserved break", C.Work.Notification.Message, "Work notification message should match")
	assert.Equal(t, `C:\User\image.svg`, C.Work.Notification.Icon, "Work notification icon should be expanded")

	// break task
	assert.Equal(t, 15*time.Minute, C.Break.Duration, "Break duration should be 15 minutes")
	assert.Equal(t, "Relaxation break", C.Break.Title, "Break title should match")
	expectedBreakThen := [][]string{{"echo", "break finished"}}
	assert.Equal(t, expectedBreakThen, C.Break.Then, "Break then commands should match")

	// break notification
	assert.False(t, C.Break.Notification.Enabled, "Break notification should be disabled")
	assert.Equal(t, "Break Over! 😴", C.Break.Notification.Title, "Break notification title should match")
	assert.Equal(t, "Back to productive work", C.Break.Notification.Message, "Break notification message should match")
	assert.Equal(t, "/abs/path/break-icon.png", C.Break.Notification.Icon, "Break notification icon should match")
}

func TestExpandPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err, "Failed to get home directory")

	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "should expand tilde at start of path",
			input: "~/my/icon.png",
			want:  homeDir + "/my/icon.png",
		},
		{
			name:  "should leave absolute unix path unchanged",
			input: "/some/file.txt",
			want:  "/some/file.txt",
		},
		{
			name:  "should leave windows absolute path unchanged",
			input: `C:\Users\user\image.svg`,
			want:  `C:\Users\user\image.svg`,
		},
		{
			name:  "should leave empty path unchanged",
			input: "",
			want:  "",
		},
		{
			name:  "should expand tilde with trailing slash to home directory",
			input: "~/",
			want:  homeDir,
		},
		{
			name:  "should not expand standalone tilde without slash",
			input: "~",
			want:  "~",
		},
		{
			name:  "should not expand tilde when part of filename",
			input: "~file",
			want:  "~file",
		},
		{
			name:  "should expand tilde in nested directory structure",
			input: "~/documents/projects/pomo/config.yml",
			want:  homeDir + "/documents/projects/pomo/config.yml",
		},
		{
			name:  "should leave relative path unchanged",
			input: "./config/file.yml",
			want:  "./config/file.yml",
		},
		{
			name:  "should not expand tilde in middle of path",
			input: "/path/~user/file.txt",
			want:  "/path/~user/file.txt",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := expandPath(tt.input)
			assert.Equal(t, tt.want, result, "Path expansion result mismatch for input: %s", tt.input)
		})
	}
}

func setupViper() {
	viper.Reset()
	C = Config{} // reset global config

	viper.SetConfigName(AppName)
	viper.SetConfigType("yaml")

	setDefaults()
}

func writeAndLoadConfig(t *testing.T, config string) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, ConfigFile)

	err := os.WriteFile(configFile, []byte(config), 0o644)
	assert.NoError(t, err, "Failed to write test config")

	viper.AddConfigPath(tempDir)
	assert.NoError(t, LoadConfig(), "Failed to load config")
}

func getDefaultConfig() Config {
	// Create a temporary viper instance to unmarshal defaults
	tempViper := viper.New()
	for key, value := range DefaultConfig {
		tempViper.SetDefault(key, value)
	}

	var config Config
	if err := tempViper.Unmarshal(&config); err != nil {
		log.Printf("Failed to unmarshal default config: %v", err)
		// Return empty config on error
		return Config{}
	}

	return config
}

// assertConfigMatches validates all fields between expected and actual config
func assertConfigMatches(t *testing.T, expected Config, actual Config) {
	// main config assertion
	assert.Equal(t, expected.AskToContinue, actual.AskToContinue)

	// ASCII Art assertions
	assert.Equal(t, expected.ASCIIArt.Enabled, actual.ASCIIArt.Enabled)
	assert.Equal(t, expected.ASCIIArt.Font, actual.ASCIIArt.Font)
	assert.Equal(t, expected.ASCIIArt.Color, actual.ASCIIArt.Color)

	// work task assertions
	assert.Equal(t, expected.Work.Duration, actual.Work.Duration)
	assert.Equal(t, expected.Work.Title, actual.Work.Title)
	assert.Equal(t, expected.Work.Then, actual.Work.Then)

	// work notification assertions
	assert.Equal(t, expected.Work.Notification.Enabled, actual.Work.Notification.Enabled)
	assert.Equal(t, expected.Work.Notification.Title, actual.Work.Notification.Title)
	assert.Equal(t, expected.Work.Notification.Message, actual.Work.Notification.Message)

	// handle path expansion for icon paths
	expectedWorkIcon, _ := expandPath(expected.Work.Notification.Icon)
	assert.Equal(t, expectedWorkIcon, actual.Work.Notification.Icon)

	// break task assertions
	assert.Equal(t, expected.Break.Duration, actual.Break.Duration)
	assert.Equal(t, expected.Break.Title, actual.Break.Title)
	assert.Equal(t, expected.Break.Then, actual.Break.Then)

	// break notification assertions
	assert.Equal(t, expected.Break.Notification.Enabled, actual.Break.Notification.Enabled)
	assert.Equal(t, expected.Break.Notification.Title, actual.Break.Notification.Title)
	assert.Equal(t, expected.Break.Notification.Message, actual.Break.Notification.Message)

	// handle path expansion for icon paths
	expectedBreakIcon, _ := expandPath(expected.Break.Notification.Icon)
	assert.Equal(t, expectedBreakIcon, actual.Break.Notification.Icon)
}
