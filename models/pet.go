package models

type Pet struct {
	ID             int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name           string    `gorm:"column:name;type:varchar(45)" json:"name"`
	Owner          string    `gorm:"column:owner;type:varchar(45)" json:"owner"`
	Species        string    `gorm:"column:species;type:varchar(45)" json:"species"`
	Sex            string    `gorm:"column:sex;type:char(1)" json:"sex"`
	Birth          Timestamp `gorm:"column:birth;type:date" json:"birth"`          // 可空日期字段
	Death          Timestamp `gorm:"column:death;type:date" json:"death"`          // 可空日期字段
	DatabaseColumn int8      `gorm:"column:database_column" json:"databaseColumn"` // 可空 bit 字段
	Del            bool      `gorm:"column:del;default:0" json:"-"`
}
