package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// UserIntent - структура распознанного намерения пользователя
type UserIntent struct {
	Action string `json:"action"` // "search" | "save" | "remind"
	Query  string `json:"query"`  // Очищенный текст запроса или заметки
	DueAt  string `json:"due_at"` // Время напоминания (если есть)
}

const systemPrompt = `You are a high-speed intent classifier for Telegram bot AICloud.
Classify user intention from voice or text messages.
Output ONLY valid JSON:
{"action":"search|save|remind","query":"extracted content in original language","due_at":""}

CLASSIFICATION RULES:

1. "search": User wants to find, recall, or check notes ("найди", "поищи", "где было", "помнишь", "покажи", "помоги найти про...").
   Keep ONLY the search subject in "query" without filler words. Keep original Russian language!
   Example: "Помоги мне найти про Дораму мой демон и сохранить." -> {"action":"search","query":"Дорама мой демон","due_at":""}
   Example: "Слушай, где-то тут был рецепт шарлотки с корицей?" -> {"action":"search","query":"рецепт шарлотки с корицей","due_at":""}

2. "remind": User wants a reminder at a specific time/date ("напомни завтра в...", "напомни вечером", "через 2 часа").
   Keep the task in "query", and time/date in "due_at". Keep original Russian language!
   Example: "Напомни завтра в 18:00 созвониться с Мишей" -> {"action":"remind","query":"созвониться с Мишей","due_at":"завтра в 18:00"}

3. "save": Notes, thoughts, shopping lists, credentials, or starts with "запиши", "сохрани", "запомни".
   Keep the core note in "query" without filler words. Keep original Russian language!
   Example: "Запиши пожалуйста код от калитки 4210" -> {"action":"save","query":"код от калитки 4210","due_at":""}
   Example: "Купить молоко, яйца и сметану" -> {"action":"save","query":"Купить молоко, яйца и сметану","due_at":""}

4. User instructions / expansions:
- If the user explicitly asks you to add details or description (e.g. "Сохрани Дорама Мой демон и добавь небольшое описание"), fulfill their instruction and include it in "query" in Russian.
`

// AnalyzerText анализирует текст через нейросеть и определяет намерение пользователя
func AnalyzerText(ctx context.Context, userText string) (UserIntent, error) {
	provider := strings.ToLower(os.Getenv("INTENT_PROVIDER"))
	if provider == "" {
		provider = "groq"
	}

	if provider == "groq" {
		intent, err := callGroqAnalyzer(ctx, userText)
		if err == nil {
			return intent, nil
		}
		log.Printf("[AnalyzerText] Ошибка Groq (%v), переключаюсь на OpenRouter fallback", err)
	}

	// По умолчанию или при откате используем OpenRouter
	return callOpenRouterAnalyzer(ctx, userText)
}

// callGroqAnalyzer отправляет запрос в Groq API (сверхбыстро, ~0.15с)
func callGroqAnalyzer(ctx context.Context, userText string) (UserIntent, error) {
	var intent UserIntent

	apiKey := os.Getenv("QROQ_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("GROQ_API_KEY")
	}
	if apiKey == "" {
		return intent, fmt.Errorf("ключ GROQ_API_KEY не задан")
	}

	model := os.Getenv("GROQ_INTENT_MODEL")
	if model == "" {
		model = "qwen/qwen3.8-27b"
	}

	reqBody, err := json.Marshal(map[string]any{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userText},
		},
		"response_format": map[string]string{"type": "json_object"},
	})
	if err != nil {
		return intent, fmt.Errorf("ошибка сериализации запроса Groq: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return intent, fmt.Errorf("ошибка создания запроса Groq: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return intent, fmt.Errorf("ошибка отправки в Groq: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return intent, fmt.Errorf("ошибка чтения ответа Groq: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return intent, fmt.Errorf("groq вернул статус %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var r struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(bodyBytes, &r); err != nil {
		return intent, fmt.Errorf("ошибка декодирования ответа Groq: %w", err)
	}

	if len(r.Choices) == 0 {
		return intent, fmt.Errorf("пустой ответ choices от Groq")
	}

	rawContent := cleanJSONContent(r.Choices[0].Message.Content)
	if err := json.Unmarshal([]byte(rawContent), &intent); err != nil {
		return intent, fmt.Errorf("не удалось распарсить JSON от Groq (%s): %w", rawContent, err)
	}

	return intent, nil
}

// callOpenRouterAnalyzer отправляет запрос в OpenRouter
func callOpenRouterAnalyzer(ctx context.Context, userText string) (UserIntent, error) {
	var intent UserIntent

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return intent, fmt.Errorf("ключ OPENROUTER_API_KEY не задан")
	}

	model := os.Getenv("INTENT_MODEL")
	if model == "" {
		model = "nvidia/nemotron-3-ultra-550b-a55b:free"
	}

	reqBody, err := json.Marshal(map[string]any{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userText},
		},
	})
	if err != nil {
		return intent, fmt.Errorf("ошибка сериализации запроса OpenRouter: %w", err)
	}

	client := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return intent, fmt.Errorf("ошибка создания запроса OpenRouter: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return intent, fmt.Errorf("ошибка отправки в OpenRouter: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return intent, fmt.Errorf("ошибка чтения ответа OpenRouter: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return intent, fmt.Errorf("openrouter вернул статус %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var r struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(bodyBytes, &r); err != nil {
		return intent, fmt.Errorf("ошибка декодирования ответа OpenRouter: %w", err)
	}

	if len(r.Choices) == 0 {
		return intent, fmt.Errorf("пустой ответ choices от OpenRouter")
	}

	rawContent := cleanJSONContent(r.Choices[0].Message.Content)
	if err := json.Unmarshal([]byte(rawContent), &intent); err != nil {
		return intent, fmt.Errorf("не удалось распарсить JSON от OpenRouter (%s): %w", rawContent, err)
	}

	return intent, nil
}

func cleanJSONContent(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}
