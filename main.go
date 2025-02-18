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

// todo: добавить paging и адаптировать код под него
// todo: сделать todo из файла api.go

func fetchRepos() (models.RepositoryList, error) {
	url := "https://api.github.com/search/repositories?q=language:Go&sort=stars&order=desc&per_page=100&page=1"

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
	var totalForks int      //Общее число форков
	var totalStars int      //Суммарное количество звезд
	var totalOpenIssues int //Всего открытых issues

	repIdList, resErr := fetchRepos()
	if resErr != nil {
		fmt.Println("Ошибка получения ответа от запроса fetchRepos:", resErr)
		return
	}

	fmt.Println("repIdList", len(repIdList.Items))

	var repositories []models.Repository

	maxRequests := 10
	semaphore := make(chan struct{}, maxRequests)
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	for i := 0; i < len(repIdList.Items); i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			// заполняем слот семафора
			semaphore <- struct{}{}
			// в конце освобождаем слот семафора
			defer func() {
				<-semaphore
			}()

			repository, err := fetchRepoById(repIdList.Items[index].ID)
			if err != nil {
				fmt.Println("ошибка при выполнении запроса fetchRepoById внутри семафора в main", err)
			}

			// записываем в слайс репозиториев репозиторий
			mutex.Lock()
			repositories = append(repositories, repository)
			mutex.Unlock()
		}(i)
	}

	wg.Wait()
	close(semaphore)

	for _, repo := range repositories {
		totalForks += repo.ForksCount
		totalStars += repo.StargazersCount
		totalOpenIssues += repo.OpenIssuesCount
	}

	fmt.Printf("Forks %d, Stars %d, OpenIssues %d \n", totalForks, totalStars, totalOpenIssues)
	fmt.Println("repositories", repositories)

}
