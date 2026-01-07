package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// updateError handles key events on the error screen
func (m model) updateError(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "enter":
		m.screen = screenAuth
		m.errorMsg = ""
		return m, nil
	}
	return m, nil
}

// viewError renders the error screen
func (m model) viewError() string {
	var b strings.Builder

	b.WriteString(errorTitleStyle.Render("Error"))
	b.WriteString("\n\n")
	b.WriteString(m.errorMsg)

	errorContent := errorBoxStyle.Render(b.String())

	// Center the error box horizontally
	errorWidth := lipgloss.Width(errorContent)
	horizontalPadding := max(0, (m.width-errorWidth)/2)

	centeredError := lipgloss.NewStyle().
		PaddingLeft(horizontalPadding).
		Render(errorContent)

	// Help footer (centered)
	helpText := "enter: back to form • ctrl+c: quit"
	help := helpStyle.Width(m.width).Align(lipgloss.Center).Render(helpText)

	// Calculate heights
	errorHeight := lipgloss.Height(centeredError)
	helpHeight := lipgloss.Height(help)

	// Create spacer to push footer to bottom
	spacerHeight := max(0, m.height-errorHeight-helpHeight)
	topPadding := spacerHeight / 2
	bottomPadding := spacerHeight - topPadding

	topSpacer := strings.Repeat("\n", topPadding)
	bottomSpacer := strings.Repeat("\n", bottomPadding)

	return topSpacer + centeredError + bottomSpacer + help
}
