package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Repository struct {
	FullName    string   `json:"full_name"`
	Description string   `json:"description"`
	Stars       int      `json:"stargazers_count"`
	Language    string   `json:"language"`
	Topics      []string `json:"topics"`
	HTMLURL     string   `json:"html_url"`
}

type searchResponse struct {
	Items []Repository `json:"items"`
}

func SearchRepositories(query string) ([]Repository, error) {
	u := fmt.Sprintf(
		"https://api.github.com/search/repositories?q=%s&sort=stars&order=desc&per_page=20",
		url.QueryEscape(query),
	)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var result searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Items, nil
}
