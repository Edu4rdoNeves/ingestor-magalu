package dependency

import (
	"github.com/Edu4rdoNeves/ingestor-magalu/infrastructure/database"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/configs/env"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	IngesterDb *gorm.DB
)

func LoadDataBases() {
	logrus.Info("Loading databases dependencies...")

	var err error

	IngesterDb, err = database.ConnectPostgre(env.IngesterDb)
	if err != nil {
		logrus.Panic(err)
	}
}
