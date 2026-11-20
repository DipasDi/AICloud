package models

type MediaInfo struct {
	TargetType string // 'photo', 'file', 'message'
	FileID     string
	FileName   string
}
