package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Repo struct {
	Name string `json:"name"`
}

func getRepos(username, token string) ([]Repo, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=200", username)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var repos []Repo
	json.Unmarshal(body, &repos)
	return repos, nil
}

func countCommits(username, repo, token string) (int, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?author=%s&per_page=100", username, repo, username)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "token "+token)
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var commits []interface{}
	json.Unmarshal(body, &commits)
	return len(commits), nil
}

func main() {
	username := os.Getenv("GITHUB_USERNAME")
	token := os.Getenv("GITHUB_TOKEN")

	if username == "" || token == "" {
		fmt.Println("Please set GITHUB_USERNAME and GITHUB_TOKEN environment variables")
		return
	}

	repos, err := getRepos(username, token)
	if err != nil {
		fmt.Println("Error fetching repos:", err)
		return
	}

	totalCommits := 0
	for _, repo := range repos {
		fmt.Printf("Fetching commits for repo: %s\n", repo.Name)
		commits, err := countCommits(username, repo.Name, token)
		if err != nil {
			fmt.Printf("Failed to fetch commits for repo %s: %v\n", repo.Name, err)
			continue
		}
		fmt.Printf("%s has %d commits\n", repo.Name, commits)
		totalCommits += commits
	}

	fmt.Printf("Total commits across all repositories: %d\n", totalCommits)
}
