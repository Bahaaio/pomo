Output pomo.gif

#Require pomo

Set Padding 20
Set Margin 20
Set BorderRadius 10

Set Theme "Catppuccin Mocha"
Set WindowBar Colorful
Set MarginFill "#0D1117"

Set FontSize 18
Set FontFamily "JetBrainsMono Nerd Font"
Set TypingSpeed 100ms

Set Framerate 60

# start work session
Type "./pomo work -t 5s"
Sleep 1s
Enter
Sleep 3s

# Adjust time up
Up 2
Sleep 2s

# Pause
Space
Sleep 2s
Space
Sleep 2s

# Adjust time down
Down 1
Sleep 2s

# Reset time
Left
Sleep 5s

# start break
Sleep 1s
Type "./pomo break"
Enter
Sleep 5s
Ctrl+c
