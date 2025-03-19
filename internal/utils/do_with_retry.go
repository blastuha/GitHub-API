package utils

import (
	"GitHubTask/internal/api"
	"fmt"
	"net/http"
	"time"
)

// - Обратить внимание на заголовки ответа GitHub API:
//    - **X-RateLimit-Remaining** – количество оставшихся запросов
//    - **Retry-After** – рекомендуемая задержка перед следующим запросом в случае превышения лимита.

// нужно отправить запрос через функцию DoWithRetry которая:
// 1. Проверяет вернул ли запрос 429 ошибку
// 2. Если ошибка есть, то запрос отправляется повторно по экспоненциальной штуке, и мы возвращаем ошибку и респонсе
// 3. Если ошибки нет, то возвращаем nil и респонс

func DoWithRetry(req *http.Request, maxRetry int) (*http.Response, error) {
	delay := 1 * time.Second
	var response *http.Response
	var respErr error

	for attempt := 1; attempt <= maxRetry; attempt++ {
		response, respErr = api.GithubClient.Do(req)

		if respErr != nil {
			time.Sleep(delay)
			delay = delay * 2
		} else {
			return response, nil
		}

		if response != nil {
			if rateLimit := response.Header.Get("X-RateLimit-Limit"); rateLimit != "" {
				fmt.Println("X-RateLimit-Limit", rateLimit)
			}
			if rateRemaning := response.Header.Get("X-RateLimit-Remaining"); rateRemaning != "" {
				fmt.Println("X-RateLimit-Remaining", rateRemaning)
			}

			if retryDelay := response.Header.Get("Retry-After"); retryDelay != "" {
				fmt.Println("retryDelay", retryDelay)
				parseDelay, parseError := time.ParseDuration(retryDelay + "s")
				if parseError != nil {
					fmt.Println("DoWithRetry, parse delay error:", parseError)
				}
				delay = parseDelay
			}
		}

	}
	return response, nil
}
