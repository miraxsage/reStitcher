package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// model is the main application model
type model struct {
	screen screen
	width  int
	height int

	// Auth form
	inputs     []textinput.Model
	focusIndex int

	// Error
	errorMsg string

	// Main screen
	list     list.Model
	viewport viewport.Model
	ready    bool
	creds    *Credentials

	// Command menu
	showCommandMenu  bool
	commandMenuIndex int
}

// NewModel creates a new application model
func NewModel() model {
	return model{
		screen:     screenAuth,
		inputs:     initAuthInputs(),
		focusIndex: 0,
	}
}

// Init initializes the model
func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, checkStoredCredentials())
}

// Update handles all messages
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle command menu if open
		if m.showCommandMenu {
			return m.updateCommandMenu(msg)
		}

		// Open command menu with "/" (except on auth screen)
		if msg.String() == "/" && m.screen != screenAuth {
			m.showCommandMenu = true
			m.commandMenuIndex = 0
			return m, nil
		}

		switch m.screen {
		case screenAuth:
			return m.updateAuth(msg)
		case screenError:
			return m.updateError(msg)
		case screenMain:
			return m.updateList(msg)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if m.screen == screenMain {
			m.updateListSize()
		}

	case checkCredsMsg:
		if msg.creds != nil {
			m.creds = msg.creds
			m.screen = screenMain
			m.initListScreen()
			m.updateListSize()
		}

	case authResultMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			m.screen = screenError
		} else {
			m.screen = screenMain
			m.initListScreen()
			m.updateListSize()
		}
	}

	// Update inputs if on auth screen (for non-KeyMsg messages like Blink)
	if m.screen == screenAuth {
		var cmd tea.Cmd
		var updatedModel tea.Model
		updatedModel, cmd = m.updateInputs(msg)
		m = updatedModel.(model)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the current screen
func (m model) View() string {
	var view string
	switch m.screen {
	case screenAuth:
		view = m.viewAuth()
	case screenError:
		view = m.viewError()
	case screenMain:
		view = m.viewList()
	}

	// Overlay command menu if open
	if m.showCommandMenu {
		view = m.overlayCommandMenu(view)
	}

	return view
}
