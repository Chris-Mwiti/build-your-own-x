package task

import (
	"time"

	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/google/uuid"
)

//here we define the states a task can be assigned
type State int

//task options
const (
	Pending State = iota
	Scheduled
	Runnig
	Completed
	Failed
)

//Task attributes and Methods structure type
type Task struct {
	ID uuid.UUID
	Name string
	State State
	Image string //represents the docker image
	Memory int //resource metric
	Disk int //resource metric
	ExposedPorts nat.PorSet //allocation of proper ports for the task
	PortBindings map[string]string //allcation of proper ports foor the task
	RestartPolicy string
	StartTime  time.Time
	Finish time.Time
}

//event that can be triggered by the task
type TaskEvent struct {
	ID uuid.UUID //unique identifier of the a task event
	State State //rep the progression of the state of th task
	Timestamp time.Time //record the time the event was requested
	Task Task //the task that requested the taskevent
}



 