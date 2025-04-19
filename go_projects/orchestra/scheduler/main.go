package scheduler

//info: Defines the contract to be implemented
type Scheduler interface {
	SelectCandidateNodes()
	Score()
	Pick()
}

