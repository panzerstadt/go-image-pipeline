package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/panzerstadt/go-image-pipeline/configs"
)

type model struct {
	choices []string
	cursor  int
}

func initialModel() model {
	return model{
		choices: []string{
			"Run Producer",
			"Run Image Processing Consumer",
			"Reset Image Processing Inputs",
			"Test Imagemagick",
			fmt.Sprintf("Create Kafka Topic: '%s'", configs.TopicImageJobs),
			fmt.Sprintf("Remove Kafka Topic: '%s'", configs.TopicImageJobs),
			"Exit"},
		cursor: 0,
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
				// runCommand("go", "run", "./cmd/producer")
				runCommand("osascript", "-e", `tell application "Terminal" to do script "cd /Users/tliqun/Documents/Github/go-image-pipeline && go run ./cmd/producer"`)
			case 1:
				runCommand("osascript", "-e", `tell application "Terminal" to do script "cd /Users/tliqun/Documents/Github/go-image-pipeline && go run ./cmd/consumer"`)
			case 2:
				inputRes := runCommand("find", "./inputs", "-type", "f")
				inputFiles := strings.Split(strings.TrimSpace(string(inputRes)), "\n")
				for _, file := range inputFiles {
					runCommand("rm", "-f", file)
				}
				outputRes := runCommand("find", "./outputs", "-type", "f")
				outputFiles := strings.Split(strings.TrimSpace(string(outputRes)), "\n")
				for _, file := range outputFiles {
					out := file
					in := strings.Replace(file, "/outputs", "/inputs", 1)
					fmt.Println(out, in)
					runCommand("cp", out, in)
					runCommand("rm", "-f", out)
				}
			case 3:
				middlePath := "./intermediate/test.jpg"
				outPath := "./outputs/test.jpg"
				runCommand("/opt/homebrew/bin/convert", "-strip", "-interlace", "Plane", "-quality", "80", "-resize", "2000x2000", middlePath, outPath)
			case 4:
				create_topic()
			case 5:
				remove_topic()
			case 6:
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

func runCommand(name string, args ...string) []byte {
	cmd := exec.Command(name, args...)
	fmt.Println(cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("%s", output)
	fmt.Print("\n\n\n\n\n")
	return output
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting app: %v\n", err)
		os.Exit(1)
	}
}
