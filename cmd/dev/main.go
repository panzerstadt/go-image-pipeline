package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/panzerstadt/go-image-pipeline/configs"
)

type logLineMsg string

type model struct {
	choices []string
	logs    []string // cmd output
	// meta
	cursor   int
	ready    bool
	viewport viewport.Model
	content  string
}

type option struct {
	label string
	value int
}

var (
	RUN_PRODUCER = option{label: "Run Producer", value: 0}
	RUN_CONSUMER = option{label: "Run Image Processing Consumer", value: 1}
	RESET_IMAGES = option{label: "Reset Image Processing Inputs", value: 2}
	TEST_DEPS    = option{label: "Test Imagemagick", value: 3}
	CREATE_TOPIC = option{label: fmt.Sprintf("Create Kafka Topic: '%s'", configs.TopicImageJobs), value: 4}
	REMOVE_TOPIC = option{label: fmt.Sprintf("Remove Kafka Topic: '%s'", configs.TopicImageJobs), value: 5}
	EXIT         = option{label: "Exit", value: 6}
)

func initialModel() model {
	return model{
		choices: []string{
			RUN_PRODUCER.label,
			RUN_CONSUMER.label,
			RESET_IMAGES.label,
			TEST_DEPS.label,
			CREATE_TOPIC.label,
			REMOVE_TOPIC.label,
			EXIT.label},
		logs: []string{},
		// metadata
		cursor:  0,
		ready:   false,
		content: "",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case logLineMsg:
		m.logs = append(m.logs, string(msg))
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
			case RUN_PRODUCER.value:
				return m, runCommandAsync("go", "run", "./cmd/producer")
				// runCommand("osascript", "-e", `tell application "Terminal" to do script "cd /Users/tliqun/Documents/Github/go-image-pipeline && go run ./cmd/producer"`)
			case RUN_CONSUMER.value:
				runCommand("osascript", "-e", `tell application "Terminal" to do script "cd /Users/tliqun/Documents/Github/go-image-pipeline && go run ./cmd/consumer"`)
			case RESET_IMAGES.value:
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
			case TEST_DEPS.value:
				middlePath := "./intermediate/test.jpg"
				outPath := "./outputs/test.jpg"
				runCommand("/opt/homebrew/bin/convert", "-strip", "-interlace", "Plane", "-quality", "80", "-resize", "2000x2000", middlePath, outPath)
			case CREATE_TOPIC.value:
				create_topic()
			case REMOVE_TOPIC.value:
				remove_topic()
			case EXIT.value:
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		choicesHeight := lipgloss.Height(m.choicesView())
		verticalMarginHeight := choicesHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = choicesHeight
			m.viewport.SetContent(m.content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	return fmt.Sprintf("%s\n%s", m.choicesView(), m.viewport.View())
}

func (m model) choicesView() string {
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

func runCommandAsync(name string, args ...string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command(name, args...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return logLineMsg(fmt.Sprintf("Error: %v", err))
		}
		if err := cmd.Start(); err != nil {
			return logLineMsg(fmt.Sprintf("Start error: %v", err))
		}

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			return logLineMsg(line)
		}

		if err := cmd.Wait(); err != nil {
			return logLineMsg(fmt.Sprintf("Command error: %v", err))
		}

		return nil
	}
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
	p := tea.NewProgram(initialModel(),
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion()) // turn on mouse support so we can track the mouse wheel

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting app: %v\n", err)
		os.Exit(1)
	}
}
