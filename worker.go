package main

import (
	"github.com/estenssoros/goflow/state"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	name string
	kill chan struct{}
}

func NewWorker() *Worker {
	return &Worker{
		name: namer.randomName(),
		kill: make(chan struct{}),
	}
}

func (w *Worker) Start(taskQueue <-chan TaskInterface, evalQueue chan<- TaskInterface) {
	logrus.Infof("starter worker %s", w.name)
	defer func() {
		logrus.Infof("worker %s exited", w.name)
	}()
	for {
		select {
		case task := <-taskQueue:
			task.SetState(state.Running)
			logrus.Infof("%s running %s", w.name, task.FormattedID())
			err := task.Run()
			if err != nil {
				logrus.Errorf("%s failed", task.FormattedID())
				task.SetState(state.Failed)
			} else {
				logrus.Infof("%s success", task.FormattedID())
				task.SetState(state.Success)
			}
			evalQueue <- task
		case <-w.kill:
			logrus.Infof("killing worker %s...", w.name)
			return
		}
	}
}
