# figma-kit demo — pure commands building the demo site

Output demo.gif

Set FontSize 20
Set FontFamily "Hack"
Set Width 1400
Set Height 800
Set Padding 24
Set Margin 24
Set MarginFill "#0D0F17"
Set BorderRadius 16
Set WindowBar Colorful
Set WindowBarSize 44
Set TypingSpeed 28ms
Set Framerate 60
Set Shell zsh
Set Theme { "name": "Noir", "black": "#0A0B12", "red": "#FF5959", "green": "#33CC80", "yellow": "#FF9933", "blue": "#3366FF", "magenta": "#8B5CF6", "cyan": "#14B8A6", "white": "#F5F7FA", "brightBlack": "#242938", "brightRed": "#FF5959", "brightGreen": "#33CC80", "brightYellow": "#FF9933", "brightBlue": "#3366FF", "brightMagenta": "#8B5CF6", "brightCyan": "#14B8A6", "brightWhite": "#F5F7FA", "background": "#0A0B12", "foreground": "#F5F7FA", "selection": "#242938", "cursor": "#3366FF" }

Require figma-kit

Hide
Type "PS1='%F{blue}❯%f '"
Enter
Sleep 300ms
Type "clear"
Enter
Sleep 300ms
Show

Sleep 400ms

Type "figma-kit --version"
Enter
Sleep 1.2s

Type "figma-kit page create 'Landing Page' -p 0"
Enter
Sleep 1.2s

Type "figma-kit page create 'Components' -p 1"
Enter
Sleep 1.2s

Type "figma-kit page create 'Design System' -p 2"
Enter
Sleep 1.2s

Type "figma-kit node create frame --name 'Hero' -w 1440 --height 900"
Enter
Sleep 1.2s

Type "figma-kit ui badge -t noir --label 'v0.2.0'"
Enter
Sleep 1.2s

Type "figma-kit ui button -t noir --variant primary --label 'Get Started'"
Enter
Sleep 1.2s

Type "figma-kit card glass -t noir --title 'AI-Native Workflow' --desc 'Prompt to production'"
Enter
Sleep 1.2s

Type "figma-kit card glass -t noir --title '120+ Commands' --desc 'Every primitive and template'"
Enter
Sleep 1.2s

Type "figma-kit card glass -t noir --title 'Terminal + AI' --desc 'One binary, any MCP agent'"
Enter
Sleep 1.2s

Type "figma-kit ds create -t noir"
Enter
Sleep 1.5s

Type "figma-kit ds variables-create -t noir --name 'Noir Tokens'"
Enter
Sleep 1.5s

Type "figma-kit export tokens --format css -t noir"
Enter
Sleep 2s

Type "figma-kit validate theme noir"
Enter
Sleep 1.5s

Sleep 3s
