package models

import "fmt"

type RepositoryList struct {
	Items []RepositoryID `json:"items"`
}

func (r *RepositoryList) PrintItems() {
	for _, item := range r.Items {
		fmt.Println("id", item.ID)
	}
}

type RepositoryID struct {
	ID int `json:"id"`
}
