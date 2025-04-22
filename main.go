package main

import (
	"flag"
	"sync"

	"github.com/Edu4rdoNeves/ingestor-magalu/cmd/worker"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/dependency"
	"github.com/Edu4rdoNeves/ingestor-magalu/utils"
	"github.com/sirupsen/logrus"
)

func main() {

	workerFlag, scriptFlag := utils.ConfigFlags()

	flag.Parse()
	wg := new(sync.WaitGroup)
	dependency.Load()

	if *workerFlag {
		wg.Add(1)
		go func() {
			worker.Run()
			wg.Done()
		}()
		wg.Wait()
	}

	if *scriptFlag {
		dependency.SimulatorTask.Run()
		logrus.Info("Finish script!")
	}
}
