package goflow

import (
	"fmt"
	"testing"
)

func TestNewDag(t *testing.T) {
	dag, err := NewDag(&DagConfig{
		ID:               "test",
		Description:      "a test to see how run",
		ScheduleInterval: "* * * * *",
	})
	if err != nil {
		t.Fatal(err)
	}

	t1, err := dag.NewBash(&BashOperator{
		TaskID:      "print date",
		BashCommand: "date",
	})
	if err != nil {
		t.Fatal(err)
	}

	t2, err := dag.NewBash(&BashOperator{
		TaskID:      "sleep",
		BashCommand: "sleep 5",
		Retries:     3,
	})
	if err != nil {
		t.Fatal(err)
	}

	t3, err := dag.NewBash(&BashOperator{
		TaskID:      "hello world",
		BashCommand: "echo hello world",
	})
	if err != nil {
		t.Fatal(err)
	}

	t4, err := dag.NewBash(&BashOperator{
		TaskID:      "hello world2",
		BashCommand: "echo hello world",
	})
	if err != nil {
		t.Fatal(err)
	}

	t5, err := dag.NewGo(&GoOperator{
		TaskID: "go program",
		GoFunc: func() error { fmt.Println("hello from go!"); return nil },
	})
	if err != nil {
		t.Fatal(err)
	}

	t2.SetUpstream(t1)
	t3.SetUpstream(t1)
	t4.SetUpstream(t3)
	t5.SetUpstream(t4)

	if err := dag.TreeView(); err != nil {
		t.Fatal(err)
	}
	scheduler := NewScheduler()

	scheduler.AddDag(dag)

	if err := scheduler.Run(); err != nil {
		t.Fatal(err)
	}
}
