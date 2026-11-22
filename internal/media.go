package internal

import (
	"AICloud/models"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/mymmrac/telego"
)

// ExtractMediaInfo извлекает тип, ID файла, имя и расширение/тип файла из сообщения Telegram
func ExtractMediaInfo(msg telego.Message) models.MediaInfo {
	switch {
	case len(msg.Photo) > 0:
		return models.MediaInfo{
			TargetType: "photo",
			FileType:   "photo",
			FileID:     msg.Photo[len(msg.Photo)-1].FileID,
		}
	case msg.Document != nil:
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(msg.Document.FileName), "."))
		if ext == "" && msg.Document.MimeType != "" {
			parts := strings.Split(msg.Document.MimeType, "/")
			if len(parts) == 2 {
				ext = strings.ToLower(parts[1])
			}
		}
		return models.MediaInfo{
			TargetType: "file",
			FileType:   ext,
			FileID:     msg.Document.FileID,
			FileName:   msg.Document.FileName,
		}
	case msg.Voice != nil:
		return models.MediaInfo{
			TargetType: "voice",
			FileType:   "voice",
			FileID:     msg.Voice.FileID,
		}
	case msg.Video != nil:
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(msg.Video.FileName), "."))
		if ext == "" {
			ext = "mp4"
		}
		return models.MediaInfo{
			TargetType: "file",
			FileType:   ext,
			FileID:     msg.Video.FileID,
			FileName:   msg.Video.FileName,
		}
	case msg.Audio != nil:
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(msg.Audio.FileName), "."))
		if ext == "" {
			ext = "mp3"
		}
		return models.MediaInfo{
			TargetType: "file",
			FileType:   ext,
			FileID:     msg.Audio.FileID,
			FileName:   msg.Audio.FileName,
		}
	default:
		return models.MediaInfo{
			TargetType: "message",
			FileType:   "text",
		}
	}
}

// GetFileBytes скачивает файл по ссылке и возвращает байты
func GetFileBytes(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка скачивания файла: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("сервер Telegram вернул статус %d при скачивании файла", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения байтов файла: %w", err)
	}
	return body, nil
}
