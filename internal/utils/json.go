package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// PrettyPrintJSON красиво форматирует JSON
func PrettyPrintJSON(body []byte) (string, error) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, body, "", "  ")
	if err != nil {
		return "", fmt.Errorf("ошибка форматирования JSON: %w", err)
	}
	return prettyJSON.String(), nil
}
