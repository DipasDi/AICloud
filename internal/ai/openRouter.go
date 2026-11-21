package ai

import (
	"AICloud/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// sendOpenRouterRequest выполняет базовый HTTP-запрос к API OpenRouter
func sendOpenRouterRequest(endpoint string, reqBody []byte) ([]byte, error) {
	client := &http.Client{
		Timeout: 80 * time.Second,
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

// GetEmbedding получает векторное представление текста через Gemini (по умолчанию) или OpenRouter
func GetEmbedding(msg string) ([]float32, error) {
	provider := strings.ToLower(os.Getenv("EMBED_PROVIDER"))
	if provider == "" || provider == "gemini" {
		if os.Getenv("GEMINI_API_KEY") != "" {
			emb, err := getGeminiEmbedding(msg)
			if err == nil {
				return emb, nil
			}
			// При ошибке Gemini, если задан OpenRouter, пробуем fallback
			if os.Getenv("OPENROUTER_API_KEY") != "" {
				return getOpenRouterEmbedding(msg)
			}
			return nil, err
		}
	}

	return getOpenRouterEmbedding(msg)
}

// getGeminiEmbedding получает вектор через Google Gemini API (1500 запросов в минуту бесплатно)
func getGeminiEmbedding(msg string) ([]float32, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY не задан")
	}

	model := os.Getenv("EMBED_MODEL")
	if model == "" || !strings.HasPrefix(model, "models/") {
		model = "models/gemini-embedding-001"
	}

	payload := map[string]any{
		"model": model,
		"content": map[string]any{
			"parts": []map[string]string{
				{"text": msg},
			},
		},
		"outputDimensionality": 768,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации запроса Gemini: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/%s:embedContent?key=%s", model, apiKey)
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса Gemini: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка отправки запроса в Gemini: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа Gemini: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini вернул статус %d: %s", resp.StatusCode, string(respBytes))
	}

	var r struct {
		Embedding struct {
			Values []float32 `json:"values"`
		} `json:"embedding"`
	}
	if err := json.Unmarshal(respBytes, &r); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа Gemini: %w", err)
	}

	if len(r.Embedding.Values) == 0 {
		return nil, fmt.Errorf("пустой массив embedding от Gemini")
	}

	return r.Embedding.Values, nil
}

// getOpenRouterEmbedding получает векторное представление текста через OpenRouter
func getOpenRouterEmbedding(msg string) ([]float32, error) {
	model := os.Getenv("EMBED_MODEL")
	if model == "" || strings.HasPrefix(model, "models/") {
		model = "nvidia/nemotron-3-embed-1b:free"
	}

	body, err := json.Marshal(models.EmbedRequest{
		Model: model,
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
