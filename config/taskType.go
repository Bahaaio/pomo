package config

type TaskType int

const (
	WorkTask TaskType = iota
	BreakTask
)

func (t TaskType) GetTask() *Task {
	if t == BreakTask {
		return &C.Break
	}
	return &C.Work
}

func (t TaskType) Opposite() TaskType {
	if t == WorkTask {
		return BreakTask
	}
	return WorkTask
}

func (t TaskType) String() string {
	return t.GetTask().Title
}
