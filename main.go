package main

import (
	"fmt"
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
	done        bool
	index       int
	width       int
	height      int
	questions   []Question
	answerField textinput.Model
	styles      *Styles
}

type Question struct {
	question string
	answer   string
}

func New(questions []Question) *model {
	answerField := textinput.New()
	answerField.Placeholder = "Enter your answer"
	answerField.Focus()

	styles := DefaultStyles()

	return &model{
		questions:   questions,
		answerField: answerField,
		index:       0,
		styles:      styles,
		done:        false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Next() {
	if m.index < len(m.questions)-1 {
		m.index++
		m.answerField.SetValue("")
		m.answerField.Focus()
	} else {
		m.index = 0
		m.answerField.SetValue("done")
		m.answerField.Blur()
	}
}

func DefaultStyles() *Styles {
	return &Styles{
		BorderColor: lipgloss.Color("#008000"),
		TextColor:   lipgloss.Color("#ffffff"),
		InputField: lipgloss.NewStyle().
			Border(lipgloss.InnerHalfBlockBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("#008000")).
			Padding(1, 2).
			Width(100),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	current := &m.questions[m.index]

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.index == len(m.questions)-1 {
				m.done = true
			}
			current.answer = m.answerField.Value()
			m.answerField.SetValue("")
			m.Next()
			return m, nil
		}
	}
	m.answerField, cmd = m.answerField.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.width == 0 {
		return "loading"
	}
	if m.done {
		var output string
		for _, question := range m.questions {
			output += fmt.Sprintf("%s: %s\n", question.question, question.answer)
		}
		return output
	}
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.questions[m.index].question,
			m.styles.InputField.Render(m.answerField.View()),
		),
	)
}

func main() {
	questions := []Question{
		{question: "What is the name of your cluster", answer: ""},
		{question: "How many control plane nodes?", answer: ""},
		{question: "How many worker nodes?", answer: ""},
	}
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
