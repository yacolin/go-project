package models

import "time"

type File struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FileName  string    `json:"file_name"`
	FilePath  string    `json:"file_path"`
	MimeType  string    `json:"mime_type"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}
