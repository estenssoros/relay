package relay

import (
	"context"

	"github.com/estenssoros/relay/state"
	"github.com/sirupsen/logrus"
)

// Worker multiprocessing unit that performs the actual task
type Worker struct {
	name string
	kill chan struct{}
}

// NewWorker creates a new worker with a clever name
func NewWorker() *Worker {
	return &Worker{
		name: namer.randomName(),
		kill: make(chan struct{}),
	}
}

// Start starts a worker workin
func (w *Worker) Start(ctx context.Context, taskQueue <-chan TaskInterface, evalQueue chan<- TaskInterface) {
	logrus.Debugf("starter worker %s", w.name)
	defer func() {
		logrus.Debugf("worker %s exited", w.name)
	}()
	for {
		select {
		case task := <-taskQueue:
			task.SetState(state.Running)
			logrus.Infof("%s running %s", w.name, task.FormattedID())
			err := task.Run()
			if err != nil {
				logrus.Errorf("%s failed", task.FormattedID())
				task.GetModel().Message = err.Error()
				task.SetState(state.Failed)
			} else {
				logrus.Infof("%s success", task.FormattedID())
				task.SetState(state.Success)
			}
			evalQueue <- task
		case <-w.kill:
			logrus.Infof("killing worker %s...", w.name)
			return
		case <-ctx.Done():
			return
		}
	}
}
