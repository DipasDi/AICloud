package models

import "time"

// Item представляет запись из таблицы items
type Item struct {
	ID             int64     `json:"id"`                        // ID записи в базе (BIGSERIAL)
	UserID         int64     `json:"user_id"`                   // ID пользователя в Telegram (BIGINT)
	ChatID         int64     `json:"chat_id"`                   // ID чата в Telegram (BIGINT)
	TargetType     string    `json:"target_type"`               // file, message, photo, link, voice
	FileType       string    `json:"file_type"`                 // exe, zip, pdf, photo, voice, text etc.
	TgFileID       string    `json:"tg_file_id"`                // ID файла в Telegram (если прикреплен)
	TgMessageID    int64     `json:"tg_message_id"`             // ID сообщения в Telegram
	FileName       string    `json:"file_name"`                 // Оригинальное имя файла
	Comment        string    `json:"comment"`                   // Текст заметки, описание или OCR-текст
	GroupName      string    `json:"group_name"`                // Название папки / группы
	CreatedAt      time.Time `json:"created_at"`                // Дата создания
	UpdatedAt      time.Time `json:"updated_at"`                // Дата последнего изменения
	Embedding      []float32 `json:"embedding,omitempty"`       // Вектор (vector(768))
	EmbeddingModel string    `json:"embedding_model,omitempty"` // Название модели эмбеддингов
	Rank           float32   `json:"rank"`                      // Релевантность
}
