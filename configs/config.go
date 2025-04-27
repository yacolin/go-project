package configs

import (
	"fmt"
	"log"

	"go-project/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectMysql() {
	cfg := &models.Config{
		User:   "root",
		Pwd:    "123456",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "springboot-learn",
	}

	dsn := fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8&parseTime=True&loc=Local", cfg.User, cfg.Pwd, cfg.Net, cfg.Addr, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected")

	DB = db
}
