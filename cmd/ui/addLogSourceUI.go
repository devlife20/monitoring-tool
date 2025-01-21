package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/devlife20/monitoring-tool/CBL"
	"github.com/devlife20/monitoring-tool/types"
	utilities "github.com/devlife20/monitoring-tool/utilies"
	"io"
	"os"
	"strings"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4).MarginTop(1)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).MarginTop(1).Foreground(lipgloss.Color("#FAFAFA"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	inputStyle        = lipgloss.NewStyle().MarginTop(1).MarginLeft(4)
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

type screenState int

const (
	selectingService screenState = iota
	selectingCloudProvider
	enteringAWSCredentials
	enteringAzureCredentials
	enteringGCPCredentials
	selectingSchedule
	finished
)

type model struct {
	list           list.Model
	cloudList      list.Model
	scheduleList   list.Model
	mainChoice     string
	cloudChoice    string
	scheduleChoice string
	quitting       bool
	state          screenState
	inputs         []textinput.Model
	credentials    types.Credentials
}

func initialModel() model {
	// Main service items
	mainItems := []list.Item{
		item("Cloud"),
		item("Local"),
		item("ELK"),
	}

	// Cloud provider items
	cloudItems := []list.Item{
		item("AWS"),
		item("Azure"),
		item("GCP"),
	}

	// Schedule items
	scheduleItems := []list.Item{
		item("Real-time"),
		item("Every 30 minutes"),
		item("Every hour"),
		item("Daily"),
	}

	const defaultWidth = 20

	mainList := list.New(mainItems, itemDelegate{}, defaultWidth, listHeight)
	mainList.Title = "Select Log Source"
	mainList.SetShowStatusBar(false)
	mainList.SetFilteringEnabled(false)
	mainList.Styles.Title = titleStyle

	cloudList := list.New(cloudItems, itemDelegate{}, defaultWidth, listHeight)
	cloudList.Title = "Select Cloud Provider"
	cloudList.SetShowStatusBar(false)
	cloudList.SetFilteringEnabled(false)
	cloudList.Styles.Title = titleStyle

	scheduleList := list.New(scheduleItems, itemDelegate{}, defaultWidth, listHeight)
	scheduleList.Title = "Select Schedule"
	scheduleList.SetShowStatusBar(false)
	scheduleList.SetFilteringEnabled(false)
	scheduleList.Styles.Title = titleStyle

	// Initialize inputs for credentials
	inputs := make([]textinput.Model, 3)
	for i := range inputs {
		t := textinput.New()
		t.CharLimit = 156
		t.Width = 20
		inputs[i] = t
	}

	return model{
		list:         mainList,
		cloudList:    cloudList,
		scheduleList: scheduleList,
		inputs:       inputs,
		state:        selectingService,
	}
}

func (m *model) setupAWSInputs() {
	m.inputs[0].Placeholder = "AWS Access Key ID"
	m.inputs[1].Placeholder = "AWS Secret Access Key"
	m.inputs[2].Placeholder = "AWS Region"
	for i := range m.inputs {
		m.inputs[i].Focus()
	}
}

func (m *model) setupAzureInputs() {
	m.inputs[0].Placeholder = "Azure Tenant ID"
	m.inputs[1].Placeholder = "Azure Client ID"
	m.inputs[2].Placeholder = "Azure Client Secret"
	for i := range m.inputs {
		m.inputs[i].Focus()
	}
}

func (m *model) setupGCPInputs() {
	m.inputs[0].Placeholder = "GCP Project ID"
	m.inputs[1].Placeholder = "Service Account Key File Path"
	m.inputs[2].Placeholder = "Region"
	for i := range m.inputs {
		m.inputs[i].Focus()
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.cloudList.SetWidth(msg.Width)
		m.scheduleList.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case selectingService:
			switch keypress := msg.String(); keypress {
			case "q", "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			case "enter":
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.mainChoice = string(i)
					if m.mainChoice == "Cloud" {
						m.state = selectingCloudProvider
					} else {
						m.state = selectingSchedule
					}
				}
				return m, nil
			}

		case selectingCloudProvider:
			switch keypress := msg.String(); keypress {
			case "esc":
				m.state = selectingService
				return m, nil
			case "enter":
				i, ok := m.cloudList.SelectedItem().(item)
				if ok {
					m.cloudChoice = string(i)
					switch m.cloudChoice {
					case "AWS":
						m.setupAWSInputs()
						m.state = enteringAWSCredentials
					case "Azure":
						m.setupAzureInputs()
						m.state = enteringAzureCredentials
					case "GCP":
						m.setupGCPInputs()
						m.state = enteringGCPCredentials
					}
				}
				return m, nil
			}

		case enteringAWSCredentials, enteringAzureCredentials, enteringGCPCredentials:
			switch keypress := msg.String(); keypress {
			case "tab":
				// Cycle through inputs
				for i := range m.inputs {
					if m.inputs[i].Focused() {
						m.inputs[i].Blur()
						nextIndex := (i + 1) % len(m.inputs)
						m.inputs[nextIndex].Focus()
						break
					}
				}
				return m, nil
			case "enter":
				// Store credentials based on provider
				switch m.state {
				case enteringAWSCredentials:
					m.credentials.AccessKeyID = m.inputs[0].Value()
					m.credentials.SecretAccessKey = m.inputs[1].Value()
					m.credentials.Region = m.inputs[2].Value()

				case enteringAzureCredentials:
					m.credentials.TenantID = m.inputs[0].Value()
					m.credentials.ClientID = m.inputs[1].Value()
					m.credentials.ClientSecret = m.inputs[2].Value()
				case enteringGCPCredentials:
					m.credentials.ProjectID = m.inputs[0].Value()
					m.credentials.KeyFilePath = m.inputs[1].Value()
					m.credentials.Region = m.inputs[2].Value()
				}
				m.state = selectingSchedule
				return m, nil
			case "esc":
				m.state = selectingCloudProvider
				return m, nil
			}

			// Handle input updates
			var cmd tea.Cmd
			for i := range m.inputs {
				if m.inputs[i].Focused() {
					m.inputs[i], cmd = m.inputs[i].Update(msg)
					break
				}
			}
			return m, cmd

		case selectingSchedule:
			switch keypress := msg.String(); keypress {
			case "esc":
				if m.mainChoice == "Cloud" {
					switch m.cloudChoice {
					case "AWS":
						m.state = enteringAWSCredentials
					case "Azure":
						m.state = enteringAzureCredentials
					case "GCP":
						m.state = enteringGCPCredentials
					}
				} else {
					m.state = selectingService
				}
				return m, nil
			case "enter":
				i, ok := m.scheduleList.SelectedItem().(item)
				if ok {
					m.scheduleChoice = string(i)

					cloudConfig := CBL.CloudLogConfigurations{
						CloudLogCredentials: m.credentials,
						Schedule:            m.scheduleChoice, // Add schedule if needed
					}
					CBL.SaveCloudConfiguration(cloudConfig, "AWS")
					m.state = finished
				}
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	switch m.state {
	case selectingService:
		m.list, cmd = m.list.Update(msg)
	case selectingCloudProvider:
		m.cloudList, cmd = m.cloudList.Update(msg)
	case selectingSchedule:
		m.scheduleList, cmd = m.scheduleList.Update(msg)
	}
	return m, cmd
}

func (m *model) View() string {
	switch m.state {
	case selectingService:
		return "\n" + m.list.View()
	case selectingCloudProvider:
		return "\n" + m.cloudList.View() + "\n(Press esc to go back)"
	case enteringAWSCredentials, enteringAzureCredentials, enteringGCPCredentials:
		var title string
		switch m.state {
		case enteringAWSCredentials:
			title = "AWS"
		case enteringAzureCredentials:
			title = "Azure"
		case enteringGCPCredentials:
			title = "GCP"
		}

		inputs := ""
		for i := range m.inputs {
			inputs += inputStyle.Render(m.inputs[i].View()) + "\n"
		}

		return fmt.Sprintf(
			"\nEnter %s Credentials:\n%s\n(Tab to switch fields, Enter to confirm, Esc to go back)",
			title,
			inputs,
		)
	case selectingSchedule:
		return "\n" + m.scheduleList.View() + "\n(Press esc to go back)"
	case finished:
		var details string
		if m.mainChoice == "Cloud" {
			details = fmt.Sprintf(
				"Selected %s on %s with schedule: %s\nCredentials have been saved.",
				m.cloudChoice,
				m.mainChoice,
				m.scheduleChoice,
			)
		} else {
			details = fmt.Sprintf(
				"Selected %s with schedule: %s",
				m.mainChoice,
				m.scheduleChoice,
			)
		}
		return quitTextStyle.Render(details)
	default:
		return "An error occurred"
	}
}

func Show() {
	m := initialModel()
	if _, err := tea.NewProgram(&m).Run(); err != nil {
		utilities.LogBubbleTeaError(err)
		os.Exit(1)
	}
}
