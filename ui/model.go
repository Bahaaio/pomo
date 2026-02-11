package ui

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Bahaaio/pomo/config"
	"github.com/Bahaaio/pomo/db"
	"github.com/Bahaaio/pomo/ui/ascii"
	"github.com/Bahaaio/pomo/ui/colors"
	"github.com/Bahaaio/pomo/ui/confirm"
	"github.com/Bahaaio/pomo/ui/summary"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	// components
	progressBar   progress.Model
	confirmDialog confirm.Model
	help          help.Model

	// timer
	timer    timer.Model
	duration time.Duration
	elapsed  time.Duration

	// state
	width, height    int // window dimensions
	onSessionEnd     string
	sessionState     SessionState
	confirmStartTime time.Time
	currentTaskType  config.TaskType
	currentTask      config.Task
	sessionSummary   summary.SessionSummary
	isShortSession   bool
	longBreak        config.LongBreak
	cyclePosition    int             // for long break tracking
	commandsWg       *sync.WaitGroup // post commands wg
	commandsCancel   context.CancelFunc

	// ASCII art
	useTimerArt     bool
	timerFont       ascii.Font
	asciiTimerStyle lipgloss.Style

	// databse
	repo *db.SessionRepo
}

func NewModel(taskType config.TaskType, cfg config.Config) Model {
	task := taskType.GetTask()

	var timerFont ascii.Font
	timerStyle := lipgloss.NewStyle()

	if cfg.ASCIIArt.Enabled {
		timerFont = ascii.GetFont(cfg.ASCIIArt.Font)

		timerColor := colors.GetColor(cfg.ASCIIArt.Color)
		timerStyle = timerStyle.Foreground(timerColor)
	}

	sessionSummary := summary.SessionSummary{}

	database, err := db.Connect()
	var repo *db.SessionRepo

	if err != nil {
		// gracefully handle database connection failure
		// fallback to in-memory summary only (nil repo)
		log.Printf("failed to initialize database: %v", err)

		// mark database as unavailable in the session summary
		sessionSummary.SetDatabaseUnavailable()
	} else {
		repo = db.NewSessionRepo(database)
	}

	return Model{
		progressBar:   progress.New(progress.WithDefaultGradient()),
		confirmDialog: confirm.New(),
		help:          help.New(),

		timer:    timer.New(task.Duration),
		duration: task.Duration,

		onSessionEnd:    cfg.OnSessionEnd,
		sessionState:    Running,
		currentTaskType: taskType,
		currentTask:     *task,
		sessionSummary:  sessionSummary,
		longBreak:       cfg.LongBreak,
		cyclePosition:   1,

		useTimerArt:     cfg.ASCIIArt.Enabled,
		timerFont:       timerFont,
		asciiTimerStyle: timerStyle,

		repo: repo,
	}
}

type SessionState byte

const (
	Running SessionState = iota
	Paused
	ShowingConfirm
	WaitingForCommands // waiting for post commands before quitting
	Quitting
)

func (m Model) GetSessionSummary() summary.SessionSummary {
	return m.sessionSummary
}
