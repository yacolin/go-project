package models

type Book struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement;"`
	ISBN        string    `json:"isbn"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Stock       int32     `json:"stock"`
	Publisher   string    `json:"publisher,omitempty"`
	PublishDate Timestamp `json:"publish_date,omitempty"`
	CreatedAt   Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt   Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
}

type BookForm struct {
	ISBN        string    `json:"isbn" binding:"required,min=5,max=15"`
	Title       string    `json:"title" binding:"required,min=2,max=50"`
	Author      string    `json:"author" binding:"required,min=2,max=50"`
	Stock       int32     `json:"stock"`
	Publisher   string    `json:"publisher,omitempty"`
	PublishDate Timestamp `json:"publish_date,omitempty"`
}

func (r *BookForm) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"isbn":         r.ISBN,
		"title":        r.Title,
		"author":       r.Author,
		"stock":        r.Stock,
		"publisher":    r.Publisher,
		"publish_date": r.PublishDate,
	}
}
