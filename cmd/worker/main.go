package worker

import (
	"context"
	"os"
	"os/signal"
	"sync"

	cronworker "github.com/Edu4rdoNeves/ingestor-magalu/cmd/worker/cron_worker"
	"github.com/Edu4rdoNeves/ingestor-magalu/cmd/worker/worker"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/configs/env"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/constants"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/dependency"
	"github.com/sirupsen/logrus"
)

func Run() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Infof("Panic recovered: %v\n", r)
		}
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	cronManager := cronworker.NewCronManager(ctx)
	cronManager.AddTask(constants.ScheduleSavePulseTask, env.ScheduleSavePulse, dependency.SavePulseTask.Run)

	workerManager := worker.NewWorkerManager(ctx, &wg)
	workerManager.AddTask(dependency.PulseTask.Run)

	cronManager.Start()
	workerManager.Start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	<-sigs
	logrus.Info("shutdown signal received, waiting for workers to complete...")

	cronManager.Stop()
	cancel()
	wg.Wait()

	logrus.Info("shutdown completed")
}
