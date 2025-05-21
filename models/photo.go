package models

type Photo struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement;"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	AlbumID     int64     `json:"album_id" gorm:"index"`
	Album       Album     `json:"-" gorm:"foreignKey:AlbumID"`
	CreatedAt   Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt   Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
}

type PhotoForm struct {
	Title       string `json:"title" binding:"required,min=2,max=100"`
	URL         string `json:"url" binding:"required,url"`
	Description string `json:"description" binding:"max=500"`
	AlbumID     int64  `json:"album_id" binding:"required"`
}

func (r *PhotoForm) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"title":       r.Title,
		"url":         r.URL,
		"description": r.Description,
		"album_id":    r.AlbumID,
	}
}
