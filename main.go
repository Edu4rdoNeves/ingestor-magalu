package main

import (
	"flag"
	"sync"

	"github.com/Edu4rdoNeves/ingestor-magalu/cmd/api/server"
	"github.com/Edu4rdoNeves/ingestor-magalu/cmd/worker"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/dependency"
	"github.com/Edu4rdoNeves/ingestor-magalu/utils"
	"github.com/sirupsen/logrus"
)

func main() {

	flags := utils.ConfigFlags()

	flag.Parse()
	wg := new(sync.WaitGroup)
	dependency.Load()

	if flags.RunWorker {
		wg.Add(1)
		go func() {
			worker.Run()
			wg.Done()
		}()
		wg.Wait()
	}

	if flags.RunAPI {
		logrus.Info("Api Run!")
		server.Run()
	}
}
