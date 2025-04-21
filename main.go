package main

import (
	"sync"

	"github.com/Edu4rdoNeves/ingestor-magalu/cmd/worker"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/configs/env"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/dependency"
	"github.com/sirupsen/logrus"
)

func main() {
	wg := new(sync.WaitGroup)

	env.LoadEnv()
	dependency.Load()

	if dependency.IngesterDb == nil {
		logrus.Panic("IngesterDb is nil after ConnectPostgre")
	}

	if env.IsWorker() {
		wg.Add(1)
		go func() {
			worker.Run()
			wg.Done()
		}()

		wg.Wait()
	}

	if env.IsScript() {
		dependency.SimulatorTask.Run()
		logrus.Info("Finish script!")
	}
}
