package main

// Screen represents the current screen state
type screen int

const (
	screenAuth screen = iota
	screenError
	screenMain
)

// Credentials stored in keyring
type Credentials struct {
	GitLabURL string `json:"gitlab_url"`
	Email     string `json:"email"`
	Token     string `json:"token"`
}

// Messages for tea.Msg
type authResultMsg struct {
	err error
}

type checkCredsMsg struct {
	creds *Credentials
}

// ListItem represents a list item for the main screen
type listItem struct {
	title, desc string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.desc }
func (i listItem) FilterValue() string { return i.title }
