package quester

import "github.com/looplab/fsm"

type Mission struct {
	Name string
	Next string
	Start func()
	End func()
}

func (m *Mission) createLeaveCallback(callbacks fsm.Callbacks) {
	if m.End != nil {
		callbacks["leave_" + m.Name] = func(e *fsm.Event) {
			m.End()
		}
	}
}

func (m *Mission) createEnterCallback(callbacks fsm.Callbacks) {
	if m.Start != nil {
		callbacks["enter_" + m.Name] = func(e *fsm.Event) {
			m.Start()
		}
	}
}
