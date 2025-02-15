package main

import (
	"GitHubTask/internal"
	"GitHubTask/internal/api"
	"GitHubTask/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

//var client = &http.Client{}

// setHeaders добавляет стандартные заголовки в HTTP-запрос
func setHeaders(req *http.Request, token string) {
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
}

func fetchRepos(token string) (models.RepositoryList, error) {
	url := "https://api.github.com/search/repositories?q=language:Go&sort=stars&order=desc"

	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return models.RepositoryList{}, fmt.Errorf("ошибка при создании запроса fetchRepos: %w", reqErr)
	}

	// Устанавливаем заголовки через общую функцию
	setHeaders(req, token)

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

func fetchRepoById(id int, token string) error {
	url := "https://api.github.com/repositories/" + strconv.Itoa(id)

	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return fmt.Errorf("ошибка при fetchRepoById: %w", reqErr)
	}

	// Используем общую функцию для установки заголовков
	setHeaders(req, token)

	resp, respErr := api.GithubClient.Do(req)
	if respErr != nil {
		return fmt.Errorf("ошибка при выполнении запроса fetchRepoById: %w", respErr)
	}
	defer resp.Body.Close()

	return nil
}

func main() {
	token, tokenErr := internal.GetGithubToken()
	if tokenErr != nil {
		fmt.Println("Ошибка получения токена в main:", tokenErr)
		return
	}

	repIdList, resErr := fetchRepos(token)
	if resErr != nil {
		fmt.Println("Ошибка получения ответа от запроса fetchRepos:", resErr)
		return
	}

	maxRequests := 10
	semaphore := make(chan struct{}, maxRequests)
	wg := sync.WaitGroup{}

	for i := 0; i < len(repIdList.Items); i++ {
		wg.Add(1)
		semaphore <- struct{}{} // Блокируем слот

		go func(repoId int) {
			defer wg.Done()
			defer func() { <-semaphore }() // Освобождаем слот

			err := fetchRepoById(repoId, token)
			if err != nil {
				fmt.Println("Ошибка при вызове fetchRepoById:", err)
			}
		}(repIdList.Items[i].ID)
	}

	wg.Wait()
}
