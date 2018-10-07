package quester

type Tasks []*Task

type Task struct {
	Statement string
	Clue string
	Resolve func(answer string) bool
}
