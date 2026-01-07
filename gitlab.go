package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// GitLabClient handles GitLab API requests
type GitLabClient struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewGitLabClient creates a new GitLab API client
func NewGitLabClient(baseURL, token string) *GitLabClient {
	return &GitLabClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		token:   token,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// GetUserEmails retrieves the authenticated user's emails
func (c *GitLabClient) GetUserEmails() ([]string, error) {
	url := c.baseURL + "/api/v4/user/emails"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("PRIVATE-TOKEN", c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return nil, fmt.Errorf("invalid token: authentication failed")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitLab API error: status %d", resp.StatusCode)
	}

	var emails []struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	result := make([]string, len(emails))
	for i, e := range emails {
		result[i] = e.Email
	}

	return result, nil
}

// ValidateCredentials checks if the credentials are valid and email matches
func ValidateCredentials(creds Credentials) error {
	client := NewGitLabClient(creds.GitLabURL, creds.Token)

	emails, err := client.GetUserEmails()
	if err != nil {
		return err
	}

	// Check if provided email matches any of the user's emails
	for _, email := range emails {
		if strings.EqualFold(email, creds.Email) {
			return nil
		}
	}

	return fmt.Errorf("email '%s' not found in your GitLab account", creds.Email)
}
