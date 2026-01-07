package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// initAuthInputs creates the text inputs for the auth form
func initAuthInputs() []textinput.Model {
	inputs := make([]textinput.Model, 3)

	// GitLab URL input
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "https://gitlab.com"
	inputs[0].Focus()
	inputs[0].CharLimit = 256
	inputs[0].Width = 40

	// Email input
	inputs[1] = textinput.New()
	inputs[1].Placeholder = "user@example.com"
	inputs[1].CharLimit = 256
	inputs[1].Width = 40

	// Token input
	inputs[2] = textinput.New()
	inputs[2].Placeholder = "glpat-xxxxxxxxxxxxxxxxxxxx"
	inputs[2].CharLimit = 256
	inputs[2].Width = 40
	inputs[2].EchoMode = textinput.EchoPassword

	return inputs
}

// checkStoredCredentials checks if credentials exist in keyring
func checkStoredCredentials() tea.Cmd {
	return func() tea.Msg {
		creds, err := LoadCredentials()
		if err != nil {
			return checkCredsMsg{creds: nil}
		}
		return checkCredsMsg{creds: creds}
	}
}

// validateCredentialsCmd validates credentials against GitLab API
func validateCredentialsCmd(creds Credentials) tea.Cmd {
	return func() tea.Msg {
		if err := ValidateCredentials(creds); err != nil {
			return authResultMsg{err: err}
		}

		// Save credentials on successful validation
		if err := SaveCredentials(creds); err != nil {
			return authResultMsg{err: err}
		}

		return authResultMsg{err: nil}
	}
}

// updateAuth handles key events on the auth screen
func (m model) updateAuth(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "tab", "down":
		m.focusIndex++
		if m.focusIndex > len(m.inputs) {
			m.focusIndex = 0
		}
		return m.updateFocus(), nil

	case "shift+tab", "up":
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs)
		}
		return m.updateFocus(), nil

	case "enter":
		if m.focusIndex == len(m.inputs) {
			// Submit button focused
			creds := Credentials{
				GitLabURL: m.inputs[0].Value(),
				Email:     m.inputs[1].Value(),
				Token:     m.inputs[2].Value(),
			}

			// Basic validation
			if creds.GitLabURL == "" || creds.Email == "" || creds.Token == "" {
				m.errorMsg = "All fields are required"
				m.screen = screenError
				return m, nil
			}

			return m, validateCredentialsCmd(creds)
		}
		// Move to next field on enter
		m.focusIndex++
		if m.focusIndex > len(m.inputs) {
			m.focusIndex = 0
		}
		return m.updateFocus(), nil
	}

	// For all other keys (character input), update the focused text input
	return m.updateInputs(msg)
}

// updateFocus updates which input has focus
func (m model) updateFocus() model {
	for i := range m.inputs {
		if i == m.focusIndex {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
	return m
}

// updateInputs updates all text inputs
func (m model) updateInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

// viewAuth renders the auth screen
func (m model) viewAuth() string {
	var b strings.Builder

	// Form title
	b.WriteString(formTitleStyle.Render("GitLab Authentication"))
	b.WriteString("\n")

	// Input fields
	labels := []string{"GitLab URL", "Email", "Personal Access Token"}
	for i, input := range m.inputs {
		b.WriteString(inputLabelStyle.Render(labels[i]))
		b.WriteString("\n")
		b.WriteString(input.View())
		b.WriteString("\n\n")
	}

	// Submit button
	submitStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230"))

	if m.focusIndex == len(m.inputs) {
		submitStyle = submitStyle.
			Background(lipgloss.Color("205")).
			Bold(true)
	}

	b.WriteString(submitStyle.Render("Submit"))

	// Wrap in form box
	formContent := formStyle.Render(b.String())

	// Center the form horizontally
	formWidth := lipgloss.Width(formContent)
	horizontalPadding := max(0, (m.width-formWidth)/2)

	centeredForm := lipgloss.NewStyle().
		PaddingLeft(horizontalPadding).
		Render(formContent)

	// Help footer (centered)
	helpText := "tab/↑↓: navigate • enter: submit/next • ctrl+c: quit"
	help := helpStyle.Width(m.width).Align(lipgloss.Center).Render(helpText)

	// Calculate heights
	formHeight := lipgloss.Height(centeredForm)
	helpHeight := lipgloss.Height(help)

	// Create spacer to push footer to bottom
	spacerHeight := max(0, m.height-formHeight-helpHeight)
	topPadding := spacerHeight / 2
	bottomPadding := spacerHeight - topPadding

	topSpacer := strings.Repeat("\n", topPadding)
	bottomSpacer := strings.Repeat("\n", bottomPadding)

	return topSpacer + centeredForm + bottomSpacer + help
}
