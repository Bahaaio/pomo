# pomo — Terminal Pomodoro Timer

![Demo](.github/assets/pomo.gif)

[![Latest Release](https://img.shields.io/github/release/Bahaaio/pomo.svg)](https://github.com/Bahaaio/pomo/releases/latest)
![Build Status](https://github.com/Bahaaio/pomo/actions/workflows/build.yml/badge.svg)

A simple, customizable Pomodoro timer for your terminal, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- 🍅 Work and break timer sessions
- 📊 Real-time progress bar visualization
- ⌨️ Keyboard shortcuts to adjust time mid-session
- ⏸️ Pause and resume timer functionality
- 🖥️ Optional full screen or inline mode
- 🔔 Cross-platform desktop notifications
- 🎨 Clean, minimal terminal UI
- 🛠️ Custom commands when timers complete

## Usage

Work sessions:

```bash
pomo              # Default work session (25m)
pomo 30m          # Custom duration
```

Break sessions:

```bash
pomo break        # Default break (5m)
pomo break 10m    # Custom duration
```

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

<details>
<summary>📁 Config file search order</summary>

pomo looks for its config file in the following order:

1. **Current directory**: `pomo.yaml` (highest priority)
2. **System config directory**:
   - **Linux**: `~/.config/pomo/pomo.yaml`
   - **macOS**: `~/Library/Application Support/pomo/pomo.yaml`
   - **Windows**: `%APPDATA%\pomo\pomo.yaml`
3. **Built-in defaults** if no config file is found

</details>

Example `pomo.yaml`:

```yaml
fullScreen: true

work:
  duration: 25m
  title: work session

  # cross-platform notifications
  notification:
    enabled: true
    title: work finished 🎉
    message: time to take a break
    icon: ~/my/icon.png

break:
  duration: 5m

  # will run after the session ends
  then:
    - [spd-say, "Back to work!"]
```

Check out [pomo.yml](pomo.yml) for a full example with all options.

### Key Bindings

| Key            | Action                    |
| -------------- | ------------------------- |
| `↑` / `k`      | Increase time by 1 minute |
| `Space`        | Pause/Resume timer        |
| `←` / `h`      | Reset to initial duration |
| `q` / `Ctrl+C` | Quit                      |

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
