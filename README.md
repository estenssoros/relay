# relay

schedule and run tasks from a single executable


## Getting started

```bash
export $RELAY_HOME=~/relay
relay initdb
```

This will start a sqlite3 database at `~/relay` and seed a config file there

## Example

```go
package main

import (
	"log"

	"github.com/estenssoros/relay"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
)

func run() error {
	dag, err := relay.NewDag(&relay.DagConfig{
		ID:               "test",
		Description:      "a test to see how run",
		ScheduleInterval: "* * * * *",
	})
	if err != nil {
		return errors.Wrap(err, "new dag")
	}

	t1, err := dag.NewBash(&relay.BashOperator{
		TaskID:      "print date",
		BashCommand: "date",
	})
	if err != nil {
		return errors.Wrap(err, "t1")
	}

	t2, err := dag.NewBash(&relay.BashOperator{
		TaskID:      "sleep",
		BashCommand: "sleep 5",
		Retries:     3,
	})
	if err != nil {
		return errors.Wrap(err, "t2")
	}

	t3, err := dag.NewBash(&relay.BashOperator{
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
	scheduler := relay.NewScheduler()

	scheduler.AddDag(dag)

	if err := scheduler.Run(); err != nil {
		return errors.Wrap(err, "scheduler run")
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
```

Tree view will show you a simple tree representation of the dag

```bash
--------------------------------------------------
DAG[test] TREE VIEW
--------------------------------------------------
print date
	sleep
	hello world
--------------------------------------------------
```

`scheduler.Run()` will start an echo webserver on port `3000` along with a scheduler and executor workers.

Dags schedules are defined using chron syntax from https://github.com/gorhill/cronexpr

## TODO

- build out web pages and html so user can interface with dag data, scheduler, etc.
- check for folder, config on startup and create if not exists
- support multiple databases for relay database
- database operators (MySQL, PostGres, MsSQL, Snowflake)
- s3 operator 
