package models

type Song struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	AlbumID     int64     `json:"album_id" gorm:"index"` // 加索引，提高查询速度
	Title       string    `json:"title"`
	Duration    int       `json:"duration"`     // 单位：秒
	TrackNumber int       `json:"track_number"` // 专辑内排序
	CreatedAt   Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt   Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
}
