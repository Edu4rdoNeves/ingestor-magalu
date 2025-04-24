package cronworker

import (
	"context"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type CronManager struct {
	ctx  context.Context
	cron *cron.Cron
}

func NewCronManager(ctx context.Context) *CronManager {
	return &CronManager{
		ctx:  ctx,
		cron: cron.New(),
	}
}

type TaskWithContext func(ctx context.Context)

func (cm *CronManager) AddTask(title, schedule string, task TaskWithContext) {
	_, err := cm.cron.AddFunc(schedule, func() {
		select {
		case <-cm.ctx.Done():
			logrus.Warnf("task %s not executed: context canceled", title)
			return
		default:
			task(cm.ctx)
		}
	})
	if err != nil {
		logrus.Errorf("Error adding cron job: %s", schedule)
	}

	logrus.Infof("task %s scheduled for: %s", title, schedule)
}

func (cm *CronManager) Start() {
	logrus.Info("Starting Cron Manager...")
	cm.cron.Start()
}

func (cm *CronManager) Stop() {
	logrus.Info("Stopping Cron Manager...")
	cm.cron.Stop()
}
