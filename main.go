package main

import (
	"GitHubTask/internal/api"
	"GitHubTask/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

func fetchRepos() (models.RepositoryList, error) {
	url := "https://api.github.com/search/repositories?q=language:Go&sort=stars&order=desc"

	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return models.RepositoryList{}, fmt.Errorf("ошибка при создании запроса fetchRepos: %w", reqErr)
	}

	// Устанавливаем заголовки через общую функцию
	api.SetHeaders(req)

	resp, respErr := api.GithubClient.Do(req)
	if respErr != nil {
		return models.RepositoryList{}, fmt.Errorf("ошибка при выполнении запроса fetchRepos: %w", respErr)
	}
	defer resp.Body.Close()

	result, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return models.RepositoryList{}, fmt.Errorf("ошибка при чтении тела запроса fetchRepos: %w", readErr)
	}

	var response models.RepositoryList
	unmarshalErr := json.Unmarshal(result, &response)
	if unmarshalErr != nil {
		return models.RepositoryList{}, fmt.Errorf("ошибка при десериализации JSON: %w", unmarshalErr)
	}

	return response, nil
}

func fetchRepoById(id int) (models.Repository, error) {
	url := "https://api.github.com/repositories/" + strconv.Itoa(id)

	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return models.Repository{}, fmt.Errorf("ошибка при fetchRepoById: %w", reqErr)
	}

	api.SetHeaders(req)

	resp, respErr := api.GithubClient.Do(req)
	if respErr != nil {
		return models.Repository{}, fmt.Errorf("ошибка при выполнении запроса fetchRepoById: %w", respErr)
	}
	defer resp.Body.Close()

	result, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return models.Repository{}, fmt.Errorf("ошибка при выполнении чтения body в fetchRepoById: %w", readErr)
	}

	var response models.Repository
	unmarhsalErr := json.Unmarshal(result, &response)
	if unmarhsalErr != nil {
		return models.Repository{}, fmt.Errorf("ошибка при десериализации JSON: %w", unmarhsalErr)
	}

	return response, nil
}

func main() {

	repIdList, resErr := fetchRepos()
	if resErr != nil {
		fmt.Println("Ошибка получения ответа от запроса fetchRepos:", resErr)
		return
	}

	//repIdList.PrintItems()

	maxRequests := 10
	semaphore := make(chan models.Repository, maxRequests)
	wg := sync.WaitGroup{}

	for i := 0; i < len(repIdList.Items); i++ {
		wg.Add(1)
		go func(index int) {

			repository, err := fetchRepoById(repIdList.Items[index].ID)
			if err != nil {
				fmt.Println("ошибка при выполнении запроса fetchRepoById внутри семафора в main", err)
			} else {
				semaphore <- repository
			}
		}(i)
	}

}
