package database

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectPostgre(config DbConfig) (*gorm.DB, error) {
	logrus.Info("Connecting to PostgreSQL database...")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d TimeZone=%s",
		config.Host, config.User, config.Password, config.DbName, config.Port, "America%2FSao_Paulo")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // loga as queries SQL
	})
	if err != nil {
		return nil, err
	}

	err = RunMigrate(db)
	if err != nil {
		return nil, err
	}

	logrus.Info("PostgreSQL database connected!")

	return db, nil
}
