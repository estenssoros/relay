# goflow

schedule and run tasks from a single executable


## Getting started

```
$ EXPORT $GOFLOW_HOME=~/goflow
$ goflow initdb
```

This will start a sqlite3 database at `~/goflow` and seed a config file there

## Example
```
package main

import (
    "log"
    "github.com/estenssoros/goflow"
    "github.com/pkg/errors"
)

func run() error {
    dag, err := goflow.NewDag(&DagConfig{
		ID:               "test",
		Description:      "a test to see how run",
		ScheduleInterval: "* * * * *",
	})
    if err != nil {
        return errors.Wrap(err, "new dag")
	}

    t1, err := dag.NewBash(&BashOperator{
		TaskID:      "print date",
		BashCommand: "date",
	})
	if err != nil {
		return errors.Wrap(err, "t1")
	}

	t2, err := dag.NewBash(&BashOperator{
		TaskID:      "sleep",
		BashCommand: "sleep 5",
		Retries:     3,
	})
	if err != nil {
		return errors.Wrap(err, "t2")
	}

	t3, err := dag.NewBash(&BashOperator{
		TaskID:      "hello world",
		BashCommand: "echo hello world",
	})
	if err != nil {
		return errors.Wrap(err, "t3")
	}

	t2.SetUpstream(t1)
	t3.SetUpstream(t1)

	if err := dag.TreeView(); err != nil {
		return errors.Wrap(err, "dag tree view")
	}
	scheduler := NewScheduler()

	scheduler.AddDag(dag)

	if err := scheduler.Run(); err != nil {
		return errors.Wrap(err, "scheduler run")
	}
}

func main() {
    if err := run(); err!=nil{
        log.Fatal(err)
    }
}
```

This will start an echo webserver on port `3000` along with a scheduler and executor workers.

Dags schedules are defined using chron syntax from https://github.com/gorhill/cronexpr


## TODO

- actually interact with the database when running dags etc.
- build out web pages and html so user can interface with dag data, scheduler, etc.
- check for folder, config on startup and create if not exists
- support multiple databases for goflow database
- database operators (MySQL, PostGres, MsSQL, Snowflake)
- s3 operator