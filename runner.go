package goflow

import (
	"context"

	"github.com/estenssoros/goflow/config"
	"github.com/estenssoros/goflow/state"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// DagRunner runs dags
type DagRunner struct {
	dagChan chan *DAG
	Error   chan error
}

// NewDagRunner creates a new dag runner
func NewDagRunner() *DagRunner {
	return &DagRunner{
		dagChan: make(chan *DAG),
		Error:   make(chan error),
	}
}

// Run waits for dags on a dag chan
func (r *DagRunner) Run(ctx context.Context) {
	for {
		select {
		case dag := <-r.dagChan:
			dagRun := dag.DagRun()
			if err := dagRun.Create(); err != nil {
				r.Error <- errors.Wrap(err, "dag run create")
			}
			if err := dag.Run(ctx, dagRun); err != nil {
				r.Error <- errors.Wrapf(err, "%s", dag.FormattedID())
			}
		case <-ctx.Done():
			logrus.Info("closing dag runner...")
			return
		}
	}

}

// RunDag  sends a dag to be run
func (r *DagRunner) RunDag(dag *DAG) {
	r.dagChan <- dag
}

// TaskRunner runs tasks in a dag
type TaskRunner struct {
	evalQueue      chan TaskInterface
	taskQueue      chan TaskInterface
	Error          chan error
	Tasks          map[string]TaskInterface
	Done           chan struct{}
	success        []TaskInterface
	failed         []TaskInterface
	upstreamFailed []TaskInterface
	workers        []*Worker
}

// NewTaskRunner creates a new task runner
func NewTaskRunner(tasks map[string]TaskInterface) *TaskRunner {
	taskMap := map[string]TaskInterface{}
	for k, v := range tasks {
		taskMap[k] = v
	}
	return &TaskRunner{
		evalQueue:      make(chan TaskInterface, len(taskMap)),
		taskQueue:      make(chan TaskInterface, len(taskMap)),
		Error:          make(chan error),
		Done:           make(chan struct{}),
		Tasks:          taskMap,
		success:        []TaskInterface{},
		failed:         []TaskInterface{},
		upstreamFailed: []TaskInterface{},
	}
}

// IsDone check to see if all tasks are accounted for
func (r *TaskRunner) IsDone() bool {
	return len(r.success)+len(r.failed)+len(r.upstreamFailed) == len(r.Tasks)
}

func (r *TaskRunner) isUpstreamFailed(task TaskInterface) bool {
	for _, t := range task.upstreamList() {
		switch t.GetState() {
		case state.Failed, state.UpstreamFailed:
			return true
		}
	}
	return false
}

func (r *TaskRunner) isUpstreamSuccess(task TaskInterface) bool {
	for _, t := range task.upstreamList() {
		if t.GetState() != state.Success {
			return false
		}
	}
	return true
}

// Evaluate evaluate tasks state and distribute to workers or lists
func (r *TaskRunner) Evaluate(ctx context.Context) {
	for {
		select {
		case task := <-r.evalQueue:
			switch task.GetState() {
			case state.Queued: // start task
				// TODO update task model to be queue
				// TODO set start date to now
				r.taskQueue <- task
				continue

			case state.Success: // add to success
				// TODO updated task model to success
				// TODO set end date to now
				r.success = append(r.success, task)

			case state.Failed: // fail downstream tasks
				// TODO updated task model to failed
				// TODO set end date to now
				for _, t := range task.downstreamList() {
					r.upstreamFailed = append(r.upstreamFailed, t)
				}
				r.failed = append(r.failed, task)

			case state.Pending: // check if runnable
				if r.isUpstreamFailed(task) {
					task.SetState(state.UpstreamFailed)
					for _, t := range task.downstreamList() {
						t.SetState(state.UpstreamFailed)
					}
					r.evalQueue <- task
					// TODO UPDATE TASK MODEL
					continue
				}
				if r.isUpstreamSuccess(task) {
					task.SetState(state.Queued)
					//TODO update model
					logrus.Infof("%s sent to workers", task.FormattedID())
					r.taskQueue <- task
					continue
				}
			case state.UpstreamFailed:
				r.upstreamFailed = append(r.upstreamFailed, task)
				continue
			}

			if r.IsDone() {
				logrus.Info("all tasks done!")
				for _, w := range r.workers {
					w.kill <- struct{}{}
				}
				logrus.Info("all workers exited")
				close(r.evalQueue)
				close(r.taskQueue)
				logrus.Info("waiting on runner done")
				r.Done <- struct{}{}
				logrus.Info("runner done!")
				return
			}
			r.evalQueue <- task

		case <-ctx.Done():
			logrus.Info("shutting down task evaluator...")
			return
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// SpawnWorkers spawn workers to hand tasks
func (r *TaskRunner) SpawnWorkers(ctx context.Context) {
	numWorkers := min(config.DefaultConfig.Webserver.Workers, len(r.Tasks))
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker()
		r.workers = append(r.workers, worker)
		go worker.Start(r.taskQueue, r.evalQueue)
	}
	logrus.Infof("starting %d workers", numWorkers)
}

// Run run the task runner
func (r *TaskRunner) Run(ctx context.Context) {

	go r.Evaluate(ctx)

	r.SpawnWorkers(ctx)

	<-r.Done

}

// FinalState calculate the final state based on the length of task lists
func (r *TaskRunner) FinalState() state.State {
	if len(r.failed) != 0 {
		return state.Failed
	}
	return state.Success
}
