package worker

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type IWorker interface {
	Start()
}

type TaskFunc func(ctx context.Context)

type WorkerManager struct {
	Context   context.Context
	WaitGroup *sync.WaitGroup
	Tasks     []TaskFunc
}

func NewWorkerManager(ctx context.Context, wg *sync.WaitGroup) *WorkerManager {
	return &WorkerManager{
		Context:   ctx,
		WaitGroup: wg,
	}
}

func (wm *WorkerManager) AddTask(task TaskFunc) {
	wm.Tasks = append(wm.Tasks, task)
}

func (wm *WorkerManager) Start(interval time.Duration) {
	logrus.Info("Starting Worker Manager...")

	wm.WaitGroup.Add(1)
	defer wm.WaitGroup.Done()

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-wm.Context.Done():
				logrus.Info("Worker Manager Stopped.")
				return
			case <-ticker.C:
				for _, task := range wm.Tasks {
					logrus.Info("Executing Task...")
					task(wm.Context)
				}
			}
		}
	}()
}
