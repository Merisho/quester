package quester

import (
	"testing"
	"fmt"
)

func TestCurrentTask(t *testing.T) {
	m := &Mission{
		Tasks: Tasks{
			{
				Statement: "task 1",
			},
			{
				Statement: "task 2",
			},
		},
	}

	if m.CurrentTask().Statement != "task 1" {
		t.Fatal("Current task must be the first task defined in collection")
	}

	m.ResolveCurrentTask("right answer")

	if m.CurrentTask().Statement != "task 2" {
		t.Fatal("Current task must be changed after successful resolution of previos one")
	}

	m.ResolveCurrentTask("right answer")

	if m.CurrentTask() != nil {
		t.Fatal("In case all task are resolved current task must return nil")
	}
}

func TestFinished(t *testing.T) {
	m := &Mission{
		Tasks: Tasks{
			{
				Statement: "task 1",
			},
		},
	}

	m.ResolveCurrentTask("test")

	if !m.IsFinished() {
		t.Fatal("Mission must be finished when all tasks are resolved")
	}
}

func TestClue(t *testing.T) {
	m := &Mission{
		Tasks: Tasks{
			{
				Statement: "task 1",
				Clue: "clue 1",
			},
		},
	}

	if m.Clue() != "clue 1" {
		t.Fatal("Clue() must return clue of the current task")
	}
}
