# pomo — Terminal Pomodoro Timer

![Demo](https://raw.githubusercontent.com/Bahaaio/pomo/main/.github/assets/pomo.gif)

[![Latest Release](https://img.shields.io/github/release/Bahaaio/pomo.svg)](https://github.com/Bahaaio/pomo/releases/latest)
![Build Status](https://github.com/Bahaaio/pomo/actions/workflows/build.yml/badge.svg)

A simple, customizable Pomodoro timer for your terminal, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- 🍅 Work and break timer sessions
- 🔗 Task chaining with user confirmation prompts
- 📊 Real-time progress bar visualization
- ⌨️ Keyboard shortcuts to adjust time mid-session
- ⏸️ Pause and resume sessions
- ⏭️ Skip to next session
- 🔔 Cross-platform desktop notifications
- 🎨 Clean, minimal terminal UI with ASCII art timer fonts
- 🛠️ Custom commands when timers complete

### Desktop Notifications

pomo sends native desktop notifications when sessions complete

<details>
<summary>🔔 View notification examples</summary>

**Linux (GNOME)**

![Linux Notification](https://raw.githubusercontent.com/Bahaaio/pomo/main/.github/assets/notification_linux.png)

**Windows**

![Windows Notification](https://raw.githubusercontent.com/Bahaaio/pomo/main/.github/assets/notification_windows.jpg)

_Note: Actual notification appearance varies by operating system and desktop environment_

</details>

## Timer Fonts

<!-- prettier-ignore -->
|              **mono12**              |                  **rebel**                   |
| :----------------------------------: | :------------------------------------------: |
| ![mono12](https://raw.githubusercontent.com/Bahaaio/pomo/main/.github/assets/mono12.png) |      ![rebel](https://raw.githubusercontent.com/Bahaaio/pomo/main/.github/assets/rebel.png)      |
|               **ansi**               |                **ansiShadow**                |
|   ![ansi](https://raw.githubusercontent.com/Bahaaio/pomo/main/.github/assets/ansi.png)   | ![ansiShadow](https://raw.githubusercontent.com/Bahaaio/pomo/main/.github/assets/ansiShadow.png) |

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
# prompt to continue after session completion
# false = exit when done
askToContinue: true

asciiArt:
  # use ASCII art for timer display
  enabled: true

  # available fonts: (mono12, rebel, ansi, ansiShadow)
  # default: mono12
  font: ansiShadow

  # color of the ASCII art timer
  # hex color or "none"
  color: "#5A56E0"

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

### Sound Notifications

You can play sounds when sessions complete by running commands in the `then` section.

```yaml
work:
  then:
    - [paplay, ~/sounds/work-done.mp3] # Linux
    # - [afplay, ~/sounds/work-done.mp3] # macOS
    # - [powershell, start, work-done.mp3] # Windows
```

### Key Bindings

#### Timer Controls

| Key            | Action                    |
| -------------- | ------------------------- |
| `↑` / `k`      | Increase time by 1 minute |
| `Space`        | Pause/Resume timer        |
| `←` / `h`      | Reset to initial duration |
| `s`            | Skip to next session      |
| `q` / `Ctrl+C` | Quit                      |

> Skip button skips directly to the next session, bypassing any prompts

#### Confirmation Dialog

| Key            | Action           |
| -------------- | ---------------- |
| `y`            | Confirm (Yes)    |
| `n`            | Cancel (No)      |
| `Tab`          | Toggle selection |
| `Enter`        | Submit choice    |
| `q` / `Ctrl+C` | Quit             |

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
