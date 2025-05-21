package models

type Song struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	AlbumID     int64     `json:"album_id" gorm:"index"` // 加索引，提高查询速度
	Album       Album     `json:"-" gorm:"foreignKey:AlbumID"`
	Title       string    `json:"title"`
	Duration    int       `json:"duration"`     // 单位：秒
	TrackNumber int       `json:"track_number"` // 专辑内排序
	CreatedAt   Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt   Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
}

type SongForm struct {
	Title       string `json:"title" binding:"required,min=1,max=100"`
	Duration    int    `json:"duration" binding:"required,gte=0"`     // 单位：秒
	TrackNumber int    `json:"track_number" binding:"required,gte=1"` // 专辑内排序
	AlbumID     int64  `json:"album_id" binding:"required"`
}

func (r *SongForm) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"title":        r.Title,
		"duration":     r.Duration,
		"track_number": r.TrackNumber,
		"album_id":     r.AlbumID,
	}
}
