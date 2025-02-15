package models

import "fmt"

type RepositoryList struct {
	Items []Repository `json:"items"`
}

func (r *RepositoryList) PrintItems() {
	for _, item := range r.Items {
		fmt.Println("id", item.ID)
	}
}

type Repository struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	StargazersCount int    `json:"stargazers_count"`
	ForksCount      int    `json:"forks_count"`
	OpenIssuesCount int    `json:"open_issues_count"`
}
