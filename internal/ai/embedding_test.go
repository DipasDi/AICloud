package ai

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestGetEmbedding(t *testing.T) {
	_ = godotenv.Load("../../.env")

	emb, err := GetEmbedding("Тестовый запрос для эмбеддинга")
	if err != nil {
		t.Fatalf("GetEmbedding failed: %v", err)
	}

	if len(emb) != 768 {
		t.Fatalf("expected 768 dimensions, got %d", len(emb))
	}

	t.Logf("Success! Got %d dimensions embedding, first value: %f", len(emb), emb[0])
}
