package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type GroqTranscriptionResponse struct {
	Text string `json:"text"`
}

func AINalysisVoice(ctx context.Context, audio []byte) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if err := writer.WriteField("model", "whisper-large-v3-turbo"); err != nil {
		return "", err
	}
	if err := writer.WriteField("language", "ru"); err != nil {
		return "", err
	}

	part, err := writer.CreateFormFile("file", "voice.ogg")
	if err != nil {
		return "", err
	}
	if _, err := part.Write(audio); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.groq.com/openai/v1/audio/transcriptions", &body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+os.Getenv("QROQ_API_KEY"))

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка отправки в Groq: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("groq вернул статус %d: %s", resp.StatusCode, string(respBytes))
	}

	var result GroqTranscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("ошибка декодирования ответа Groq: %w", err)
	}

	return result.Text, nil
}
