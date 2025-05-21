package models

type Comment struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement;"`
	PhotoID   int64     `json:"photo_id" gorm:"index"`
	Photo     Photo     `json:"-" gorm:"foreignKey:PhotoID"`
	Content   string    `json:"content" gorm:"type:varchar(500);not null"`
	Author    string    `json:"author" gorm:"type:varchar(100);not null"`
	CreatedAt Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
}

type CommentForm struct {
	PhotoID int64  `json:"photo_id" binding:"required"`
	Content string `json:"content" binding:"required,min=1,max=500"`
	Author  string `json:"author" binding:"required,min=1,max=100"`
}

func (r *CommentForm) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"photo_id": r.PhotoID,
		"content":  r.Content,
		"author":   r.Author,
	}
}
