package ai

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestAnalyzerText(t *testing.T) {
	_ = godotenv.Load("../../.env")

	// 1. Проверяем работу через Groq
	os.Setenv("INTENT_PROVIDER", "groq")
	intentGroq, err := AnalyzerText(context.Background(), "Помоги мне найти про Дораму мой демон и сохранить.")
	if err != nil {
		t.Fatalf("AnalyzerText (Groq) failed: %v", err)
	}
	if intentGroq.Action != "search" {
		t.Errorf("Groq expected action 'search', got '%s'", intentGroq.Action)
	}
	t.Logf("Groq Result: action=%s, query=%s", intentGroq.Action, intentGroq.Query)
}
