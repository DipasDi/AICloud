package ai

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func TestAnalyzerText(t *testing.T) {
	_ = godotenv.Load("../../.env")
	os.Setenv("INTENT_PROVIDER", "groq")

	// 1. Поиск
	searchIntent, err := AnalyzerText(context.Background(), "Помоги мне найти про Дораму мой демон и сохранить.")
	if err != nil {
		t.Fatalf("Search test failed: %v", err)
	}
	if searchIntent.Action != "search" {
		t.Errorf("expected action 'search', got '%s'", searchIntent.Action)
	}
	t.Logf("Search result: %s", searchIntent.Query)

	// 2. Дополнение описания по запросу
	enrichIntent, err := AnalyzerText(context.Background(), "Сохрани, посмотри Дораму, мой демон С небольшим описанием Чтоб там было описание, сам найди и сам допиши")
	if err != nil {
		t.Fatalf("Enrich test failed: %v", err)
	}
	if enrichIntent.Action != "save" {
		t.Errorf("expected action 'save', got '%s'", enrichIntent.Action)
	}
	// Проверяем, что нейросеть сгенерировала описание, а не просто повторила фразу "с описанием"
	if len(enrichIntent.Query) < 40 || !strings.Contains(strings.ToLower(enrichIntent.Query), "демон") {
		t.Errorf("expected detailed enriched query, got '%s'", enrichIntent.Query)
	}
	t.Logf("Enriched result:\n%s", enrichIntent.Query)
}
