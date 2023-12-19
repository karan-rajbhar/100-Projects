package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	notes    []string         //notes
	cursor   int              //cursor position
	selected map[int]struct{} //selected notes
}

func (m *model) View() string {
	s := "Select a note to view:\n\n"
	for i, choice := range m.notes {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	return s
}

func (m *model) Init() tea.Cmd {
	return nil
}

func initialModel() *model {
	return &model{
		notes:    []string{"First", "Second", "Third"},
		selected: make(map[int]struct{}),
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case "down", "j":
			if m.cursor < len(m.notes)-1 {
				m.cursor++
			}
			return m, nil
		case "enter", " ":
			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
			return m, nil
		}
	}
	return m, nil
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
