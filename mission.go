package quester

import "github.com/looplab/fsm"

type Mission struct {
	Name string
	Next string
	Start func()
	End func()
	Tasks Tasks
	currTask int
	finished bool
}

func (m *Mission) ResolveCurrentTask(answer string) bool {
	curr := m.CurrentTask()
	res := true

	if curr.Resolve != nil {
		res = curr.Resolve(answer)
	}

	if res {
		m.nextTask()
	}

	return res
}

func (m *Mission) CurrentTask() *Task {
	l := len(m.Tasks)
	if m == nil || l == 0 || m.IsFinished() {
		return nil
	}

	return m.Tasks[m.currTask]
}
func (m *Mission) Clue() string {
	curr := m.CurrentTask()
	if curr == nil {
		return ""
	}

	return curr.Clue
}

func (m *Mission) IsFinished() bool {
	return m.finished
}

func (m *Mission) nextTask() {
	m.currTask++
	if m.currTask >= len(m.Tasks) {
		m.finished = true
	}
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
