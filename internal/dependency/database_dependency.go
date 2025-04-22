package dependency

import (
	"github.com/Edu4rdoNeves/ingestor-magalu/infrastructure/database"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/configs/env"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	IngestorDb *gorm.DB
)

func LoadDataBases() {
	logrus.Info("Loading databases dependencies...")

	var err error

	IngestorDb, err = database.ConnectPostgre(env.IngestorDb)
	if err != nil {
		logrus.Panic(err)
	}
}
