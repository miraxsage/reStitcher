package main

import "github.com/charmbracelet/lipgloss"

var (
	// Common styles
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	// Main screen styles
	sidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)

	contentStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)

	// Auth form styles
	formStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)

	formTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("62")).
			MarginBottom(1)

	inputLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	// Error styles
	errorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("9")).
			Foreground(lipgloss.Color("9")).
			Padding(1, 2)

	errorTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9"))

	// Command menu styles
	commandMenuStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62")).
				Padding(1, 2)

	commandMenuTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("62")).
				MarginBottom(1)

	commandItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	commandItemSelectedStyle = lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("205"))

	commandDescStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				MarginLeft(2)

	// Settings modal styles
	settingsTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("62"))

	settingsTabActiveStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("255")).
				Background(lipgloss.Color("62")).
				Padding(0, 2)

	settingsTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Padding(0, 2)

	settingsLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("255"))

	settingsDescStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245"))

	settingsErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("9"))

	settingsButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255")).
				Background(lipgloss.Color("240")).
				Padding(0, 2)

	settingsButtonActiveStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("255")).
					Background(lipgloss.Color("62")).
					Bold(true).
					Padding(0, 2)

	settingsButtonDisabledStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("245")).
					Background(lipgloss.Color("238")).
					Padding(0, 2)
)
