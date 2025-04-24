package worker

import (
	"context"
	"sync"

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

func (wm *WorkerManager) Start() {
	logrus.Info("Starting Worker Manager...")

	for _, task := range wm.Tasks {
		wm.WaitGroup.Add(1)

		go func(t func(ctx context.Context)) {
			defer wm.WaitGroup.Done()
			logrus.Info("Starting Task...")
			t(wm.Context)
		}(task)
	}
}
