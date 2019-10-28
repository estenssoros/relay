package goflow

import (
	"fmt"
	"io/ioutil"

	"github.com/estenssoros/goflow/models"
	"github.com/estenssoros/goflow/state"
	"github.com/pkg/errors"
)

// MySQLOperator operator for mysql
type MySQLOperator struct {
	TaskID            string
	DAG               *DAG `json:"-"` // avoid recursion
	ConnectionID      string
	Retries           int
	Message           string
	SQLCommand        string
	SQLFileLoc        string
	upstreamTaskIDs   []string
	downstreamTaskIDs []string
	State             state.State
	model             *models.TaskInstance
}

func (o *MySQLOperator) String() string { return o.TaskID }

// GetID returns the tag id for an operator
func (o *MySQLOperator) GetID() string { return o.TaskID }

// FormattedID exports the formatted id for an operator
func (o *MySQLOperator) FormattedID() string { return fmt.Sprintf("[TASK] %s", o.TaskID) }

func (o *MySQLOperator) hasUpstream() bool { return len(o.upstreamTaskIDs) > 0 }

func (o *MySQLOperator) downstreamList() []TaskInterface { return downstreamList(o) }

func (o *MySQLOperator) upstreamList() []TaskInterface { return upstreamList(o) }

// GetDag returns the dag for an operator
func (o *MySQLOperator) GetDag() *DAG { return o.DAG }

// SetDag sets the dag on an operator
func (o *MySQLOperator) SetDag(dag *DAG) { o.DAG = dag }

// HasDag checks to see if the operators dag is nil
func (o *MySQLOperator) HasDag() bool { return o.DAG != nil }

func (o *MySQLOperator) addDownstreamTask(taskID string) {
	o.downstreamTaskIDs = append(o.downstreamTaskIDs, taskID)
}

func (o *MySQLOperator) addUpstreamTask(taskID string) {
	o.upstreamTaskIDs = append(o.upstreamTaskIDs, taskID)
}

// SetState sets the state on an operator
func (o *MySQLOperator) SetState(s state.State) { o.State = s }

// GetState gets the state from an operator
func (o *MySQLOperator) GetState() state.State { return o.State }

// IsRoot checks to see if an operator has upstream tasks
func (o *MySQLOperator) IsRoot() bool { return !o.hasUpstream() }

func (o *MySQLOperator) check() error {
	if o.ConnectionID == "" {
		return errors.New("operator missing connection id")
	}
	if o.SQLCommand == "" && o.SQLFileLoc == "" {
		return errors.New("operator needs sql command or file location")
	}
	return nil
}

// Run run the bash operator
func (o *MySQLOperator) Run() error {
	if err := o.check(); err != nil {
		return errors.Wrap(err, "mysql operator check")
	}
	var sql string
	if o.SQLCommand == "" {
		b, err := ioutil.ReadFile(o.SQLFileLoc)
		if err != nil {
			return errors.Wrap(err, "read file")
		}
		sql = string(b)
	} else {
		sql = o.SQLFileLoc
	}
	conn, err := GetConnection(o.ConnectionID)
	if err != nil {
		return errors.Wrapf(err, "get connection: %s", o.ConnectionID)
	}
	defer conn.Close()
	if err := conn.Exec(sql); err != nil {
		return errors.Wrap(err, "exec sql")
	}
	return nil
}

func (o *MySQLOperator) OperatorType() string { return `mysql` }

func (o *MySQLOperator) SetModel(m *models.TaskInstance) {
	o.model = m
}

func (o *MySQLOperator) GetModel() *models.TaskInstance { return o.model }
