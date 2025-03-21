package main

import (
	"GitHubTask/internal/api"
	"GitHubTask/internal/models"
	"GitHubTask/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

// todo: добавить paging и адаптировать код под него
// todo: сделать todo из файла api.go

func fetchRepos(page int) (models.RepositoryList, error) {
	//url := "https://api.github.com/search/repositories?q=language:Go&sort=stars&order=desc&per_page=100&page=1"
	url := fmt.Sprintf("https://api.github.com/search/repositories?q=language:Go&sort=stars&order=desc&per_page=100&page=%d", page)

	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return models.RepositoryList{}, fmt.Errorf("ошибка при создании запроса fetchRepos: %w", reqErr)
	}

	// Устанавливаем заголовки через общую функцию
	api.SetHeaders(req)

	resp, respErr := utils.DoWithRetry(req, 10)
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
		return models.Repository{}, fmt.Errorf("ошибка при создании NewRequest в fetchRepoById: %w", reqErr)
	}

	api.SetHeaders(req)

	resp, respErr := utils.DoWithRetry(req, 5)

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
	var reposPagesSlice []models.RepositoryList
	reposPageNumber := 2

	for i := 1; i <= reposPageNumber; i++ {
		res, resErr := fetchRepos(i)
		if resErr != nil {
			fmt.Println("Ошибка получения ответа от запроса fetchRepos:", resErr)
			return
		}
		reposPagesSlice = append(reposPagesSlice, res)
	}

	for _, items := range reposPagesSlice {
		fmt.Println("items \n\n\n", items)
	}

	var repositories []models.Repository

	maxRequests := 10
	semaphore := make(chan struct{}, maxRequests)
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	for _, page := range reposPagesSlice {
		pageCopy := page // копируем, чтобы каждая горутина работала с уникальными данными
		for i := 0; i < len(pageCopy.Items); i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				// заполняем слот семафора
				semaphore <- struct{}{}
				// в конце освобождаем слот семафора
				defer func() {
					<-semaphore
				}()

				repository, err := fetchRepoById(pageCopy.Items[index].ID)
				if err != nil {
					fmt.Println("ошибка при выполнении запроса fetchRepoById внутри семафора в main", err)
				}

				// записываем в слайс репозиториев репозиторий
				mutex.Lock()
				repositories = append(repositories, repository)
				mutex.Unlock()
			}(i)
		}
	}

	wg.Wait()
	close(semaphore)

	for _, repo := range repositories {
		totalForks += repo.ForksCount
		totalStars += repo.StargazersCount
		totalOpenIssues += repo.OpenIssuesCount
	}

	fmt.Printf("Forks %d, Stars %d, OpenIssues %d \n", totalForks, totalStars, totalOpenIssues)
	fmt.Println("rep len", len(repositories))

}
