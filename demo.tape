Output demo.gif

Set FontSize 14
Set Width 1200
Set Height 600
Set Theme { "name": "Noir", "black": "#0D0F17", "red": "#FF5959", "green": "#33CC80", "yellow": "#FF9933", "blue": "#3366FF", "magenta": "#8B5CF6", "cyan": "#14B8A6", "white": "#F5F7FA", "brightBlack": "#242938", "brightRed": "#FF5959", "brightGreen": "#33CC80", "brightYellow": "#FF9933", "brightBlue": "#3366FF", "brightMagenta": "#8B5CF6", "brightCyan": "#14B8A6", "brightWhite": "#F5F7FA", "background": "#0D0F17", "foreground": "#F5F7FA", "selection": "#242938", "cursor": "#3366FF" }
Set Padding 20
Set Shell zsh
Set TypingSpeed 40ms

Require figma-kit

Sleep 500ms

Type "# figma-kit — Build Figma designs from the terminal"
Enter
Sleep 1s

Type "figma-kit --version"
Enter
Sleep 1.5s

Type "# Create a page"
Enter
Sleep 500ms
Type "figma-kit page create 'Landing Page' -p 0"
Enter
Sleep 2s

Type "# Generate a UI button"
Enter
Sleep 500ms
Type "figma-kit ui button -t noir --variant primary --label 'Get Started'"
Enter
Sleep 2s

Type "# Build a glass card"
Enter
Sleep 500ms
Type "figma-kit card glass -t noir --title 'Feature' --desc 'Fast and efficient'"
Enter
Sleep 2s

Type "# Create a full design system page"
Enter
Sleep 500ms
Type "figma-kit ds create -t noir"
Enter
Sleep 2s

Type "# Export theme as CSS variables"
Enter
Sleep 500ms
Type "figma-kit export tokens --format css -t noir"
Enter
Sleep 2s

Type "# Validate theme"
Enter
Sleep 500ms
Type "figma-kit validate theme noir"
Enter
Sleep 2s

Type "# That's figma-kit. Install: brew install dop-amine/tap/figma-kit"
Enter
Sleep 3s
