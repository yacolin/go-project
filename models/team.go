package models

type Team struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement;"`
	Champions int8      `json:"champions"`
	City      string    `json:"city"`
	Divide    string    `json:"divide"`
	Logo      string    `json:"logo"`
	Name      string    `json:"name"`
	Part      string    `json:"part"`
	CreatedAt Timestamp `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt Timestamp `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
}

type TeamForm struct {
	Name string `json:"name" binding:"required"`
	City string `json:"description" binding:"required"`
}

// ToMap converts TeamForm to a map for updates.
func (tf *TeamForm) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name": tf.Name,
		"city": tf.City,
	}
}
