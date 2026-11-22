package models

type MediaInfo struct {
	TargetType string // 'photo', 'file', 'message', 'voice'
	FileType   string // 'exe', 'zip', 'pdf', 'photo', 'voice', etc.
	FileID     string
	FileName   string
}
