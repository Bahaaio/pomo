# pomo — Terminal Pomodoro Timer

pomo is a simple, customizable Pomodoro timer for your terminal, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)

## Features

- Work/break timers with configurable durations
- Progress bar visualization
- Keyboard shortcuts to adjust time mid-session
- Optional alt-screen mode
- Custom commands to run when a timer ends (e.g. desktop notifications)

## Build

```bash
git clone https://github.com/Bahaaio/pomo
cd pomo
go build .
```

## Configuration

pomo looks for its config file in:

1. `~/.config/pomo/pomo.yaml`
2. Current directory (`./pomo.yaml`)

Example `pomo.yaml`:

```yaml
altScreen: true

work:
  duration: 25m
  then:
    - notify-send "Work Finished!" "Time to take a break ☕"

break:
  duration: 5m
  then:
    - notify-send "Break Over"
```

The `then` field is a list of shell commands run when the timer finishes.v

## Usage

```bash
./pomo
```

### Key Bindings

| Key    | Action                    |
| ------ | ------------------------- |
| `up`   | Increase time by 1 minute |
| `down` | Decrease time by 1 minute |
| `left` | Reset to initial duration |
| `q`    | Quit                      |

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
