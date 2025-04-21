package dependency

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var (
	CronSchedule *cron.Cron
)

func LoadWorkerDependencies() {
	logrus.Info("Loading worker dependencies...")

	tz, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		logrus.Panic("failed to load location")
	}

	CronSchedule = cron.New(cron.WithLocation(tz))
}
