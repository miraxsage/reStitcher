package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// initListScreen initializes the main list screen
func (m *model) initListScreen() {
	// Create list items
	items := []list.Item{
		listItem{title: "Introduction", desc: "Getting started guide"},
		listItem{title: "Installation", desc: "How to install the app"},
		listItem{title: "Configuration", desc: "Setup and configure"},
		listItem{title: "Usage", desc: "Basic usage examples"},
		listItem{title: "Advanced", desc: "Advanced features"},
		listItem{title: "Plugins", desc: "Available plugins"},
		listItem{title: "Themes", desc: "Customize appearance"},
		listItem{title: "API Reference", desc: "API documentation"},
		listItem{title: "FAQ", desc: "Frequently asked questions"},
		listItem{title: "Changelog", desc: "Version history"},
	}

	// Create list
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Open MRs"
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	m.list = l
	m.ready = false
}

// updateListSize updates the list and viewport dimensions
func (m *model) updateListSize() {
	if m.width == 0 || m.height == 0 {
		return
	}

	sidebarWidth := m.width / 3
	contentWidth := m.width - sidebarWidth - 4

	m.list.SetSize(sidebarWidth-4, m.height-5)

	if !m.ready {
		m.viewport = viewport.New(contentWidth-4, m.height-5)
		m.viewport.SetContent(m.renderMarkdown())
		m.ready = true
	} else {
		m.viewport.Width = contentWidth - 4
		m.viewport.Height = m.height - 5
	}
}

// updateList handles key events on the main list screen
func (m model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg.String() {
	case "ctrl+c", "q", "esc":
		return m, tea.Quit
	}

	// Handle list updates
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	// Update content when selection changes
	if m.ready {
		m.viewport.SetContent(m.renderMarkdown())
	}

	// Handle viewport updates
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// renderMarkdown renders the markdown content for the selected item
func (m model) renderMarkdown() string {
	selected := m.list.SelectedItem()
	if selected == nil {
		return ""
	}

	i := selected.(listItem)

	markdown := fmt.Sprintf(`# %s

%s

## Overview

This is the **%s** section of reStitcher documentation.

### Features

- Feature one with *emphasis*
- Feature two with **strong emphasis**
- Feature three with `+"`code`"+`

### Example

`+"```go"+`
package main

func main() {
    fmt.Println("Hello, reStitcher!")
}
`+"```"+`

> This is a blockquote with some helpful information
> about the current section.

---

*Last updated: 2025*
`, i.title, i.desc, i.title)

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(m.viewport.Width),
	)

	rendered, err := renderer.Render(markdown)
	if err != nil {
		return markdown
	}

	return rendered
}

// viewList renders the main list screen
func (m model) viewList() string {
	if !m.ready {
		return "Initializing..."
	}

	sidebarWidth := m.width / 3
	contentWidth := m.width - sidebarWidth - 4

	// Render sidebar
	sidebar := sidebarStyle.
		Width(sidebarWidth).
		Height(m.height - 4).
		Render(m.list.View())

	// Render content
	content := contentStyle.
		Width(contentWidth).
		Height(m.height - 4).
		Render(m.viewport.View())

	// Combine sidebar and content
	main := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content)

	// Help footer (centered)
	helpText := "/: commands • q/esc: quit • ↑/↓: navigate • scroll: pgup/pgdn"
	help := helpStyle.Width(m.width).Align(lipgloss.Center).Render(helpText)

	return lipgloss.JoinVertical(lipgloss.Left, main, help)
}
