package relay

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/estenssoros/relay/models"
	"github.com/estenssoros/relay/state"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// BashOperator runs a bash command
type BashOperator struct {
	TaskID            string
	DAG               *DAG `json:"-"` // avoid recursion
	Retries           int
	Message           string
	BashCommand       string
	Dir               string
	upstreamTaskIDs   []string
	downstreamTaskIDs []string
	State             state.State
	model             *models.TaskInstance
}

func (o BashOperator) String() string {
	return o.TaskID
}

// GetID returns the tag id for an operator
func (o *BashOperator) GetID() string {
	return o.TaskID
}

// FormattedID exports the formatted id for an operator
func (o *BashOperator) FormattedID() string {
	return fmt.Sprintf("[TASK] %s", o.TaskID)
}

// GetDag returns the dag for an operator
func (o *BashOperator) GetDag() *DAG {
	return o.DAG
}

// HasDag checks to see if the operators dag is nil
func (o *BashOperator) HasDag() bool {
	return o.DAG != nil
}

// SetDag sets the dag on an operator
func (o *BashOperator) SetDag(dag *DAG) {
	o.DAG = dag
}

// addDownstreamTask adds a task id to the downstream list
func (o *BashOperator) addDownstreamTask(taskID string) {
	o.downstreamTaskIDs = append(o.downstreamTaskIDs, taskID)
}

// addUpstreamTask adds a task to the upstream list
func (o *BashOperator) addUpstreamTask(taskID string) {
	o.upstreamTaskIDs = append(o.upstreamTaskIDs, taskID)
}

// SetUpstream creates relationships between tasks
func (o *BashOperator) SetUpstream(task TaskInterface) {
	setRelatives(o, task, true)
}

// SetDownStream creates relationships between tasks
func (o *BashOperator) SetDownStream(task TaskInterface) {
	setRelatives(o, task, false)
}

// hasUpstream returns true if the operators has upstream tasks
func (o *BashOperator) hasUpstream() bool {
	return len(o.upstreamTaskIDs) > 0
}

// downstreamList returns the list of downstream tasks
func (o *BashOperator) downstreamList() []TaskInterface {
	lst := []TaskInterface{}
	for _, taskID := range o.downstreamTaskIDs {
		task, err := o.DAG.getTask(taskID)
		if err != nil {
			continue
		}
		lst = append(lst, task)
	}
	return lst
}

func (o *BashOperator) upstreamList() []TaskInterface {
	lst := []TaskInterface{}
	for _, taskID := range o.upstreamTaskIDs {
		task, err := o.DAG.getTask(taskID)
		if err != nil {
			continue
		}
		lst = append(lst, task)
	}
	return lst
}

// IsRoot checks to see if an operator has upstream tasks
func (o *BashOperator) IsRoot() bool {
	return !o.hasUpstream()
}

// SetState sets the state on an operator
func (o *BashOperator) SetState(s state.State) {
	o.State = s
}

// GetState gets the state from an operator
func (o *BashOperator) GetState() state.State {
	return o.State
}

// Run run the bash operator
func (o *BashOperator) Run() error {
	cmdArgs := strings.Fields(o.BashCommand)
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

func (o *BashOperator) OperatorType() string { return `bash` }

func (o *BashOperator) SetModel(m *models.TaskInstance) {
	o.model = m
}

func (o *BashOperator) GetModel() *models.TaskInstance { return o.model }
