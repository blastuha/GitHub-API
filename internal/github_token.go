package internal

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func GetGithubToken() (string, error) {
	var githubToken string

	// Ищет файл .env в текущей директории.
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки .env файла")
	}
	// Присваиваем токен
	githubToken = os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return "", fmt.Errorf("oшибка: переменная окружения GITHUB_TOKEN не найдена")
	}

	return githubToken, nil
}
