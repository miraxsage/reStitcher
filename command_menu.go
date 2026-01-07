package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// commandItem represents a command in the menu
type commandItem struct {
	name string
	desc string
}

// commands is the list of available commands
var commands = []commandItem{
	{name: "logout", desc: "Clear your current gitlab credentials to auth again"},
}

// updateCommandMenu handles key events when command menu is open
func (m model) updateCommandMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.showCommandMenu = false
		return m, nil

	case "up", "k":
		if m.commandMenuIndex > 0 {
			m.commandMenuIndex--
		}
		return m, nil

	case "down", "j":
		if m.commandMenuIndex < len(commands)-1 {
			m.commandMenuIndex++
		}
		return m, nil

	case "enter":
		return m.executeCommand(commands[m.commandMenuIndex].name)

	case "ctrl+c":
		return m, tea.Quit
	}

	return m, nil
}

// executeCommand executes the selected command
func (m model) executeCommand(name string) (tea.Model, tea.Cmd) {
	switch name {
	case "logout":
		// Delete credentials from keyring
		DeleteCredentials()

		// Reset to auth screen
		m.showCommandMenu = false
		m.screen = screenAuth
		m.inputs = initAuthInputs()
		m.focusIndex = 0
		m.creds = nil
		m.ready = false

		return m, nil
	}

	return m, nil
}

// overlayCommandMenu renders the command menu as an overlay on top of the current view
func (m model) overlayCommandMenu(background string) string {
	// Build menu content
	var b strings.Builder

	b.WriteString(commandMenuTitleStyle.Render("Commands"))
	b.WriteString("\n")

	for i, cmd := range commands {
		var itemStyle lipgloss.Style
		prefix := "  "
		if i == m.commandMenuIndex {
			itemStyle = commandItemSelectedStyle
			prefix = "> "
		} else {
			itemStyle = commandItemStyle
		}

		b.WriteString(itemStyle.Render(prefix+cmd.name) + commandDescStyle.Render(" - "+cmd.desc))
		b.WriteString("\n")
	}

	// Help footer
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓: navigate • enter: select • q/esc: close"))

	menuContent := commandMenuStyle.Render(b.String())

	// Overlay menu on top of background (centered)
	return placeOverlayCenter(menuContent, background, m.width, m.height)
}
