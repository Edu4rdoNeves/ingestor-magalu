package database

import "gorm.io/gorm"

type DbConfig struct {
	Host     string
	Port     int64
	User     string
	Password string
	DbName   string
	Timezone string
}

var db *gorm.DB

func Get() *gorm.DB {
	return db
}
