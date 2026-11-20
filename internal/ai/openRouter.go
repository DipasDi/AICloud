package ai

import (
	"AICloud/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// sendOpenRouterRequest выполняет базовый HTTP-запрос к API OpenRouter
func sendOpenRouterRequest(endpoint string, reqBody []byte) ([]byte, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1"+endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENROUTER_API_KEY"))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка отправки POST в openrouter: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openrouter вернул статус %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// GetEmbedding получает векторное представление текста
func GetEmbedding(msg string) ([]float32, error) {
	body, err := json.Marshal(models.EmbedRequest{
		Model: os.Getenv("EMBED_MODEL"),
		Input: msg,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации запроса: %w", err)
	}

	respBytes, err := sendOpenRouterRequest("/embeddings", body)
	if err != nil {
		return nil, err
	}

	var res models.EmbedResponse
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа от openrouter: %w", err)
	}

	if len(res.Data) == 0 {
		return nil, fmt.Errorf("пустой ответ эмбеддинга от openrouter")
	}

	return res.Data[0].Embedding, nil
}

// ChatCompletion отправляет запрос к модели для обычного чата
func ChatCompletion(reqBody []byte) (string, error) {
	respBytes, err := sendOpenRouterRequest("/chat/completions", reqBody)
	if err != nil {
		return "", err
	}

	var res models.OpenRouterResponse
	if err := json.Unmarshal(respBytes, &res); err != nil {
		return "", fmt.Errorf("ошибка декодирования ответа от openrouter: %w", err)
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("пустой массив choices в ответе от openrouter")
	}

	return res.Choices[0].Message.Content, nil
}
