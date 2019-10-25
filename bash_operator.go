package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/estenssoros/goflow/state"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type BashOperator struct {
	TaskID            string
	Dag               *Dag `json:"-"` // avoid recursion
	Retries           int
	Message           string
	BashCommand       string
	upstreamTaskIDs   []string
	downstreamTaskIDs []string
	State             state.State
}

// GetID returns the tag id for an operator
func (b *BashOperator) GetID() string {
	return b.TaskID
}

// GetDag returns the dag for an operator
func (b *BashOperator) GetDag() *Dag {
	return b.Dag
}

// HasDag checks to see if the operators dag is nil
func (b *BashOperator) HasDag() bool {
	return b.Dag != nil
}

// SetDag sets the dag on an operator
func (b *BashOperator) SetDag(dag *Dag) {
	b.Dag = dag
}

// addDownstreamTask adds a task id to the downstream list
func (b *BashOperator) addDownstreamTask(taskID string) {
	b.downstreamTaskIDs = append(b.downstreamTaskIDs, taskID)
}

// addUpstreamTask adds a task to the upstream list
func (b *BashOperator) addUpstreamTask(taskID string) {
	b.upstreamTaskIDs = append(b.upstreamTaskIDs, taskID)
}

// SetUpstream creates relationships between tasks
func (b *BashOperator) SetUpstream(task TaskInterface) {
	setRelatives(b, task, true)
}

// SetDownStream creates relationships between tasks
func (b *BashOperator) SetDownStream(task TaskInterface) {
	setRelatives(b, task, false)
}

func (b BashOperator) String() string {
	return b.TaskID
}

func (b *BashOperator) FormattedID() string {
	return fmt.Sprintf("[TASK] %s", b.TaskID)
}

// hasUpstream returns true if the operators has upstream tasks
func (b *BashOperator) hasUpstream() bool {
	return len(b.upstreamTaskIDs) > 0
}

// downstreamList returns the list of downstream tasks
func (b *BashOperator) downstreamList() []TaskInterface {
	lst := []TaskInterface{}
	for _, taskID := range b.downstreamTaskIDs {
		task, err := b.Dag.getTask(taskID)
		if err != nil {
			continue
		}
		lst = append(lst, task)
	}
	return lst
}

func (b *BashOperator) upstreamList() []TaskInterface {
	lst := []TaskInterface{}
	for _, taskID := range b.upstreamTaskIDs {
		task, err := b.Dag.getTask(taskID)
		if err != nil {
			continue
		}
		lst = append(lst, task)
	}
	return lst
}

func (b *BashOperator) IsRoot() bool {
	return !b.hasUpstream()
}

// Run run the bash operator
func (b *BashOperator) Run() error {
	cmdArgs := strings.Fields(b.BashCommand)
	if len(cmdArgs) == 0 {
		return errors.New("bash command had no args")
	}
	var cmd *exec.Cmd
	if len(cmdArgs) == 1 {
		cmd = exec.Command(cmdArgs[0])
	} else {
		cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
	}

	var stderr bytes.Buffer
	mw := io.MultiWriter(&stderr, os.Stderr)
	cmd.Stderr = mw
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("%s\n%s", err, stderr.String())
	}
	logrus.Infof("running: %s (PID: %d)", strings.Join(cmd.Args, " "), cmd.Process.Pid)
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("%s\n%s", err, stderr.String())
	}
	return nil
}

func (b *BashOperator) SetState(s state.State) {
	b.State = s
}

func (b *BashOperator) GetState() state.State {
	return b.State
}
