package api

import (
	"GitHubTask/internal"
	"fmt"
	"net/http"
)

// todo: Переделать на customTransport
// todo: Есть в https://chatgpt.com/c/67adc7e2-9334-8004-beb2-f41e7bbbc7c3

var GithubClient = &http.Client{
	Transport: &http.Transport{DisableKeepAlives: false},
}

// SetHeaders добавляет стандартные заголовки в HTTP-запрос
func SetHeaders(req *http.Request) {
	token, tokenErr := internal.GetGithubToken()
	if tokenErr != nil {
		fmt.Println("Ошибка получения токена в main:", tokenErr)
		return
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
}
