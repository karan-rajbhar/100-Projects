package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type item struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (i item) GetTitle() string       { return i.Title }
func (i item) GetDescription() string { return i.Description }
func (i item) FilterValue() string    { return i.Title }

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type model struct {
	list         list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
	notesPath    string
}

func newModel() model {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	notesPath := "notes.json"

	// Check if the file exists
	if _, err := os.Stat(notesPath); os.IsNotExist(err) {
		// The file doesn't exist, create it
		_, err := os.Create(notesPath)
		if err != nil {
			log.Fatal(err)
		}

		// Create a default note
		defaultNote := item{
			Title:       "Default Note",
			Description: "This is a default note.",
		}
		// Write the default note to the file
		defaultNoteJson, err := json.Marshal([]item{defaultNote})
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile("notes.json", defaultNoteJson, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Read notes from the JSON file
	file, err := os.ReadFile(notesPath)
	if err != nil {
		log.Fatal(err)
	}

	var items []item
	err = json.Unmarshal(file, &items)
	if err != nil {
		log.Fatal(err)
	}

	var listItems []list.Item
	for _, i := range items {
		listItems = append(listItems, list.Item(i)) // Replace list.Item(i) with the correct conversion if necessary
	}

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	// Use listItems instead of items
	notesList := list.New(listItems, delegate, 0, 0)

	notesList.Title = "Notes"
	notesList.Styles.Title = titleStyle
	notesList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	return model{
		list:         notesList,
		keys:         listKeys,
		delegateKeys: delegateKeys,
		notesPath:    notesPath,
	}
}
func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		case key.Matches(msg, m.keys.insertItem):
			// Create a new note
			newNote := item{
				Title:       "New Note",
				Description: "This is a new note.",
			}

			// Add the new note to the list
			insCmd := m.list.InsertItem(0, newNote)
			statusCmd := m.list.NewStatusMessage(statusMessageStyle("Added " + newNote.GetTitle()))

			// Read the existing notes
			file, err := os.ReadFile(m.notesPath)
			if err != nil {
				log.Fatal(err)
			}

			var notes []item
			err = json.Unmarshal(file, &notes)
			if err != nil {
				log.Fatal(err)
			}

			// Add the new note to the notes
			notes = append(notes, newNote)

			// fmt.Printf("%+v\n", notes)
			notesJson, err := json.Marshal(notes)
			if err != nil {
				log.Fatal(err)
			}

			err = os.WriteFile("notes.json", notesJson, 0644)
			if err != nil {
				log.Fatal(err)
			}

			return m, tea.Batch(insCmd, statusCmd)
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.list.View())
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
