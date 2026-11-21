package internal

import (
	"AICloud/models"
	"fmt"
	"io"
	"net/http"

	"github.com/mymmrac/telego"
)

// ExtractMediaInfo извлекает тип, ID файла и имя из сообщения Telegram
func ExtractMediaInfo(msg telego.Message) models.MediaInfo {
	switch {
	case len(msg.Photo) > 0:
		return models.MediaInfo{
			TargetType: "photo",
			FileID:     msg.Photo[len(msg.Photo)-1].FileID,
		}
	case msg.Document != nil:
		return models.MediaInfo{
			TargetType: "file",
			FileID:     msg.Document.FileID,
			FileName:   msg.Document.FileName,
		}
	case msg.Voice != nil:
		return models.MediaInfo{
			TargetType: "voice",
			FileID:     msg.Voice.FileID,
		}
	default:
		return models.MediaInfo{
			TargetType: "message",
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
