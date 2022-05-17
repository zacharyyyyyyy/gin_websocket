package model

type Taskqueue struct {
	Id         int
	Type       string
	Param      string
	FailMsg    string
	Status     int
	RetryTimes int
	CreateTime int
	BeginTime  int
}

const (
	StatusNotBegin = iota
	StatusRunning
)

const (
	TypeMq = "mq"
)
