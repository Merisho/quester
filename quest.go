package quester

import (
	"errors"
	"github.com/looplab/fsm"
)

const (
	passMissionEvent = "pass"
)

var (
	NoMissionsErr = errors.New("no missions")
	QuestStartedErr = errors.New("quest is already started")
	QuestNotStartedErr = errors.New("quest not started")
	QuestFinishedErr = errors.New("quest finished")
)

func NewQuest() *Quest {
	q := &Quest{
		missions: make(map[string]*Mission),
	}
	return q
}

type Quest struct {
	initMission *Mission
	missions map[string]*Mission
	started bool
	finished bool
	fsm *fsm.FSM
}

func (q *Quest) Start() error {
	if q.isStarted() {
		return QuestStartedErr
	}
	if q.initMission == nil {
		return NoMissionsErr
	}

	q.started = true
	q.fsm = q.createFSM()

	curr := q.Current()

	if curr.Start != nil {
		curr.Start()
	}

	return nil
}

func (q *Quest) AddMission(m Mission) *Quest {
	q.missions[m.Name] = &m
	if len(q.missions) == 1 {
		q.initMission = q.missions[m.Name]
	}

	return q
}

func (q *Quest) Length() int {
	return len(q.missions)
}

func (q *Quest) Current() *Mission {
	if !q.isStarted() {
		return nil
	}

	return q.missions[q.fsm.Current()]
}

func (q *Quest) PassCurrent() error {
	if !q.isStarted() {
		return QuestNotStartedErr
	}
	if q.finished {
		return QuestFinishedErr
	}

	q.fsm.Event(passMissionEvent)

	avail := q.fsm.AvailableTransitions()
	if len(avail) == 0 {
		q.finished = true
		return nil
	}

	return nil
}

func (q *Quest) createFSM() *fsm.FSM {
	events := fsm.Events{}
	callbacks := fsm.Callbacks{}
	for _, m := range q.missions {
		if m.Next == "" {
			continue
		}

		e := fsm.EventDesc{
			Name: passMissionEvent,
			Src: []string{m.Name},
			Dst: m.Next,
		}

		events = append(events, e)
	}

	for _, m := range q.missions {
		m.createEnterCallback(callbacks)
		m.createLeaveCallback(callbacks)
	}

	return fsm.NewFSM(q.initMission.Name, events, callbacks)
}

func (q *Quest) IsFinished() bool {
	return q.finished
}

func (q *Quest) isStarted() bool {
	return q.started
}

func (q *Quest) ForEachMission(c func(m *Mission)) *Quest {
	for _, m := range q.missions {
		c(m)
	}
	return q
}
