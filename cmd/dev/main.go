package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices []string
	cursor  int
}

func initialModel() model {
	return model{
		choices: []string{"Run Producer", "Run Consumer", "Test Imagemagick", "Exit"},
		cursor:  0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			switch m.cursor {
			case 0:
				runCommand("go", "run", "./cmd/producer")
			case 1:
				runCommand("osascript", "-e", `tell application "Terminal" to do script "cd /Users/tliqun/Documents/Github/go-image-pipeline && go run ./cmd/consumer"`)
			case 2:
				middlePath := "./intermediate/test.jpg"
				outPath := "./outputs/test.jpg"
				runCommand("/opt/homebrew/bin/convert", "-strip", "-interlace", "Plane", "-quality", "80", "-resize", "2000x2000", middlePath, outPath)
			case 3:
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "What do you want to do?\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	s += "\nUse the arrow keys to navigate and press Enter to select.\n"
	return s
}

func runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("%s", output)
	fmt.Print("\n\n\n\n\n")
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting app: %v\n", err)
		os.Exit(1)
	}
}
