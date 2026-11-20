package models

import "time"

type Item struct {
	ID          int       `json:"id"`            // id записи в бд
	UserID      int       `json:"user_id"`       // id пользователя в телеграме
	ChatID      int       `json:"chat_id"`       // id чата в телеграме
	TargetType  string    `json:"target_type"`   // file, screenshot, video, note
	TgFileID    string    `json:"tg_file_id"`    // пусто, если это просто текстовая заметка
	TgMessageID int       `json:"tg_message_id"` // id сообщения в телеграме
	FileName    string    `json:"file_name"`     // оригинальное имя файла (если есть)
	Comment     string    `json:"comment"`       // текст заметки, описание или распознанный текст с картинки (OCR)
	GroupName   string    `json:"group_name"`    // группа/папка (например: "сайты")
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EmbedRequest struct {
	Model      string `json:"model"`
	Input      string `json:"input"`
	Dimensions int    `json:"dimensions,omitempty"`
}

type EmbedResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}
