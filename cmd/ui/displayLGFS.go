package ui

import (
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
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).Background(lipgloss.Color("5")).Bold(true)
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
	confirmState
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
	state       state
	list        list.Model
	logFiles    []linux.LogFile
	content     string
	showConfirm bool
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

		case "ctrl+s":
			m.state = confirmState
			return m, nil

		case "y":
			if m.state == confirmState {
				m.showConfirm = false
				m.state = listState
				fmt.Println("LFS source added successfully!")
			}
			return m, nil

		case "n":
			if m.state == confirmState {
				m.showConfirm = false
				m.state = listState
			}
			return m, nil

		case "enter":
			if m.state == listState {
				selectedIndex := m.list.Index()
				m.state = viewState
				m.content = linux.ReadFileContent(m.logFiles[selectedIndex].Path)
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
	case confirmState:
		return lipgloss.NewStyle().Margin(2).Padding(2).Bold(true).Render(
			"Are you sure you want to add LFS source? (y/n)",
		)
	}
	return ""
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
