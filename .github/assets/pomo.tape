Output pomo.gif

Require pomo

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

Sleep 1s

# start short session
Type "pomo 5s"
Sleep 0.7s
Enter
Sleep 3s

# pause
Space
Sleep 2s

# resume
Space
Sleep 2s

# show confirmation dialog interaction
Sleep 2s

# select "no" option
Right
Sleep 1s

# select "yes" option
Left
Sleep 1s
Type "y"

# show break session
Sleep 5s
Ctrl+c

# clear terminal
Hide
Type "clear"
Enter
Show
Sleep 1s

# start work session
Type pomo
Sleep 0.7s
Enter
Sleep 2s

# adjust time up
Up
Sleep 0.2s
Up
Sleep 0.2s
Up
Sleep 0.2s
Up
Sleep 0.2s
Up
Sleep 0.2s
Sleep 3s

Ctrl+c
