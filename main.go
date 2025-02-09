package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func fetchRepos() []byte {
	url := "https://api.github.com/search/repositories?q=language:Go&sort=stars&order=desc"
	// создаем свой клиент, через него мы будем вызывать req(запрос), через метод client.Do
	// здесь мы не настраиваем свой кастомный клиент, его можно не создавать, а обращаться через http.DefaultClient.Do(req)
	client := &http.Client{}

	// Загружаем токен из переменной окружения
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Ошибка загрузки .env файла")
	}
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("Ошибка: переменная окружения GITHUB_TOKEN не найдена")
		return nil
	}

	// создаем тело запроса, третьим аргументом передается body (тело запроса), данные которые мы хотим удалить или запостить через методы POST, DELETE
	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		fmt.Println("Ошибка при создании запроса fetchRepos", reqErr)
	}

	// устанавливаем заголовки по документации GitHub API для успешного запроса
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// делаем сам запрос
	resp, respErr := client.Do(req)
	if respErr != nil {
		fmt.Println("Ошибка при выполнении запроса fetchRepos", respErr)
	}
	defer resp.Body.Close() // Закрываем тело ответа, чтобы избежать утечек памяти

	// Читаем тело ответа
	result, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		fmt.Println("Ошибка при чтении тела запроса в fetchRepos", readErr)
	}

	return result
}

func main() {
	fmt.Println(string(fetchRepos()))
}
