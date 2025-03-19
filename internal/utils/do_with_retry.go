package utils

import (
	"GitHubTask/internal/api"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// - Обратить внимание на заголовки ответа GitHub API:
//    - **X-RateLimit-Remaining** – количество оставшихся запросов
//    - **Retry-After** – рекомендуемая задержка перед следующим запросом в случае превышения лимита.

func DoWithRetry(req *http.Request, maxRetry int) (*http.Response, error) {
	delay := 1 * time.Second
	var response *http.Response
	var respErr error

	for attempt := 1; attempt <= maxRetry; attempt++ {
		response, respErr = api.GithubClient.Do(req)
		// Если сетевая ошибка
		if respErr != nil {
			// Если попытка последняя
			if attempt == maxRetry {
				return nil, fmt.Errorf("последняя попытка исчерпана: %v", respErr)
			}

			time.Sleep(delay)
			delay *= 2
			continue
		}

		// Если сервер вернул 429
		if response.StatusCode == http.StatusTooManyRequests {
			if rateLimit := response.Header.Get("X-RateLimit-Limit"); rateLimit != "" {
				fmt.Println("X-RateLimit-Limit:", rateLimit)
			}
			if rateRemaining := response.Header.Get("X-RateLimit-Remaining"); rateRemaining != "" {
				fmt.Println("X-RateLimit-Remaining:", rateRemaining)
			}

			if retryDelay := response.Header.Get("Retry-After"); retryDelay != "" {
				fmt.Println("Retry-After:", retryDelay)
				if seconds, err := strconv.Atoi(retryDelay); err == nil {
					delay = time.Duration(seconds) * time.Second
				} else {
					fmt.Println("DoWithRetry, parse delay error:", err)
				}
			}
			time.Sleep(delay)
			continue
		}

		// Если ошибок нет – возвращаем ответ
		return response, nil
	}

	return nil, fmt.Errorf("все попытки исчерпаны")
}
