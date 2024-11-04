package models

type StageStatus int

type StageStatusKey struct {
	ReleaseID string
	StageName string
	ChildID   string
}

const (
	Pending          StageStatus = iota //
	InProgress                          // the child is notified
	ShouldEnd                           // only WaitForSignal stage type. The child poll for it on /end_stage to finish a stage
	WaitingForResult                    // either after ShouldEnd or after InProgress (may stay at InProgress and jump to Completed)
	Completed                           // received the stage result
	Failure                             // received the stage result as Failure
	Error                               // received the stage result as Error
)
