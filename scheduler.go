package goflow

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/estenssoros/goflow/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Scheduler orchestrates the scheduling of dags
type Scheduler struct {
	Dags       map[string]*DAG
	DagNextRun map[string]time.Time
	mu         sync.Mutex
}

// NewScheduler creates a new scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		Dags: map[string]*DAG{},
	}
}

// AddDag adds a dag to the scheduler
func (s *Scheduler) AddDag(dag *DAG) error {
	_, ok := s.Dags[dag.ID]
	if ok {
		return errors.Errorf("dag: %s allread registered", dag.ID)
	}
	if err := dag.getOrCreateDagModel(); err != nil {
		return err
	}
	s.Dags[dag.ID] = dag
	return nil
}

func (s *Scheduler) hearbeat(ctx context.Context, ch chan<- string) {
	ticker := time.NewTicker(time.Duration(config.DefaultConfig.SchedulerHeartBeatSec) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			logrus.Infof("scheduler heartbeat")
			now := time.Now().UTC()
			for dagID, nextRun := range s.DagNextRun {
				if now.Equal(nextRun) || now.After(nextRun) {
					s.Dags[dagID].updateNextScheduled()
					ch <- dagID
				}
			}
		case <-ctx.Done():
			logrus.Infof("closing heartbeat...")
			return
		}
	}
}

func (s *Scheduler) setDagNextRun() error {
	setDagNextRun := map[string]time.Time{}
	for dagID, dag := range s.Dags {
		nextRun, err := dag.NextRun()
		if err != nil {
			return errors.Wrap(err, "dag next run")
		}
		setDagNextRun[dagID] = nextRun
	}
	s.DagNextRun = setDagNextRun
	return nil
}

func (s *Scheduler) updateDagNextRun(dagID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	dag, ok := s.Dags[dagID]
	if !ok {
		return errors.Errorf("could not find dag: %s", dagID)
	}
	nextRun, err := dag.NextRun()
	if err != nil {
		return errors.Wrap(err, "dag next run")
	}
	s.DagNextRun[dagID] = nextRun
	return nil
}

// Run runs the dags in a scheduler
func (s *Scheduler) Run() error {
	webServer := NewWebserver(s.Dags)
	go webServer.Serve()

	logrus.Infof("starting scheduler heartbeat: %d seconds", config.DefaultConfig.Scheduler.SchedulerHeartBeatSec)
	if err := s.setDagNextRun(); err != nil {
		return errors.Wrap(err, "set dag time map")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dagChan := make(chan string)

	go s.hearbeat(ctx, dagChan)

	dagRunner := NewDagRunner()

	go dagRunner.Run(ctx)

	go func() {
		<-ctx.Done()
		close(dagChan)
		close(dagRunner.Error)
	}()

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)

	for {
		select {
		case dagID := <-dagChan:
			if err := s.updateDagNextRun(dagID); err != nil {
				return errors.Wrap(err, "update next dag run")
			}
			dag, ok := s.Dags[dagID]
			if !ok {
				return errors.Errorf("could not locate dag: %s", dagID)
			}
			dagRunner.RunDag(dag)
		case err := <-dagRunner.Error:
			if err != nil {
				logrus.Error(err)
			}
		case <-killSignal:
			logrus.Infof("kill signal recieved. exiting...")
			return nil
		}
	}
}
