package models

type Album struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement;"`
	Name        string    `json:"name"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Liked       int64     `json:"liked"`
	CreatedAt   Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt   Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
}

type AlbumForm struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Author      string `json:"author" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"max=500"`
	Liked       int64  `json:"liked" binding:"gte=0"`
}

type UserForm struct {
	Phone string `json:"phone" binding:"required,regexp=^182-[0-9]{4}-[0-9]{4}$"`
}

func (r *AlbumForm) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":        r.Name,
		"author":      r.Author,
		"description": r.Description,
		"liked":       r.Liked,
	}
}
