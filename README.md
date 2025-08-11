# pomo — Terminal Pomodoro Timer

![Demo](.github/assets/pomo.gif)

A simple, customizable Pomodoro timer for your terminal, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- 🍅 Work and break timer sessions
- 📊 Real-time progress bar visualization
- ⌨️ Keyboard shortcuts to adjust time mid-session
- 🖥️ Optional full screen or inline mode
- 🔔 Custom commands when timers complete (notifications, etc.)
- 🎨 Clean, minimal terminal UI

## Installation

Install with Go

```bash
go install github.com/Bahaaio/pomo@latest
```

Or, build from source

```bash
git clone https://github.com/Bahaaio/pomo
cd pomo
go build .
```

Alternatively, download pre-built binaries from the [releases page](https://github.com/Bahaaio/pomo/releases).

## Configuration

pomo looks for its config file in:

1. `~/.config/pomo/pomo.yaml`
2. Current directory (`./pomo.yaml`)

Example `pomo.yaml`:

```yaml
fullScreen: true

work:
  duration: 25m
  then:
    - notify-send "Work Finished!" "Time to take a break ☕"

break:
  duration: 5m
  then:
    - notify-send "Break Over" "Back to work! 💪"
```

The `then` field contains shell commands that run when the timer finishes.

## Usage

```bash
# Start a work session (default)
./pomo

# Explicit work session
./pomo work

# Start a break session
./pomo break
```

### Key Bindings

| Key            | Action                    |
| -------------- | ------------------------- |
| `↑` / `k`      | Increase time by 1 minute |
| `↓` / `j`      | Decrease time by 1 minute |
| `←` / `h`      | Reset to initial duration |
| `q` / `Ctrl+C` | Quit                      |

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
