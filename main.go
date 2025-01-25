package main

import (
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	BorderColor lipgloss.Color
	TextColor   lipgloss.Color
	InputField  lipgloss.Style
}

type model struct {
	questions   []string
	width       int
	height      int
	index       int
	answerField textinput.Model
	styles      *Styles
}

func New(questions []string) *model {
	answerField := textinput.New()
	answerField.Placeholder = "Enter your answer"
	answerField.Focus()

	styles := DefaultStyles()

	return &model{
		questions:   questions,
		answerField: answerField,
		index:       0,
		styles:      styles,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func DefaultStyles() *Styles {
	return &Styles{
		BorderColor: lipgloss.Color("#000000"),
		TextColor:   lipgloss.Color("#ffffff"),
		InputField: lipgloss.NewStyle().
			Border(lipgloss.InnerHalfBlockBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("#000000")).
			Padding(1, 2).
			Width(100),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "loading"
	}
	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.questions[m.index],
		m.styles.InputField.Render(m.answerField.View()),
	)
}

func main() {
	questions := []string{"What is the name of your cluster", "How many control plane nodes?", "How many worker nodes?"}
	m := New(questions)

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
