package ui

import (
	"bufio"
	"fmt"
	"github.com/devlife20/monitoring-tool/LFS/linux"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// Replace with the actual import path of your `linux` package
)

const (
	listHeight = 20
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	logContentStyle   = lipgloss.NewStyle().Margin(2).Padding(2)
)

type state int

const (
	listState state = iota
	viewState
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	state    state
	list     list.Model
	logFiles []linux.LogFile
	content  string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			if m.state == listState {
				selectedIndex := m.list.Index()
				m.state = viewState
				m.content = readFileContent(m.logFiles[selectedIndex].Path)
			}
			return m, nil

		case "backspace":
			if m.state == viewState {
				m.state = listState
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	if m.state == listState {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case listState:
		return "\n" + m.list.View()
	case viewState:
		return logContentStyle.Render(m.content)
	}
	return ""
}

func readFileContent(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Sprintf("Failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	var builder strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		builder.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		builder.WriteString(fmt.Sprintf("\nError reading file: %v", err))
	}

	return builder.String()
}

// Run initializes and starts the Bubble Tea program.
func Run(logDir string) {
	logFiles, err := linux.FetchLogFiles(logDir)
	if err != nil {
		fmt.Printf("Error fetching log files: %v\n", err)
		return
	}

	if len(logFiles) == 0 {
		fmt.Println("No log files found.")
		return
	}

	items := make([]list.Item, len(logFiles))
	for i, logFile := range logFiles {
		items[i] = item(logFile.Name)
	}

	l := list.New(items, itemDelegate{}, 0, listHeight)
	l.Title = "Available Log Files"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{
		state:    listState,
		list:     l,
		logFiles: logFiles,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
