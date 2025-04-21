package database

import (
	"github.com/Edu4rdoNeves/ingestor-magalu/domain/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetEntities() []interface{} {
	return []interface{}{
		&entity.PulseData{},
	}
}

func RunMigrate(db *gorm.DB) error {
	entities := GetEntities()

	if err := db.AutoMigrate(entities...); err != nil {
		logrus.Error("Fail to run migrations: ", err)
		return err
	}

	logrus.Info("Database migrated successfully.")
	return nil
}
