package models

type File struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FileName  string    `json:"file_name"`
	FilePath  string    `json:"file_path"`
	OssURL    string    `json:"oss_url"` // OSS访问URL
	MimeType  string    `json:"mime_type"`
	Size      int64     `json:"size"`
	CreatedAt Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
}

type FileForm struct {
	FileName string `form:"file" binding:"required"`
}

func (f *File) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"file_name": f.FileName,
		"file_path": f.FilePath,
		"oss_url":   f.OssURL,
		"mime_type": f.MimeType,
		"size":      f.Size,
	}
}
