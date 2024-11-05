package models

import "fmt"

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

var stageStatusLabels = []string{
	"Pending",
	"InProgress",
	"ShouldEnd",
	"WaitingForResult",
	"Completed",
	"Failure",
	"Error",
}

// String returns the string representation of the StageStatus (you can print as %s)
func (s StageStatus) String() string {
	if s < 0 || int(s) >= len(stageStatusLabels) {
		return fmt.Sprintf("StageStatus(%d)", s)
	}
	return stageStatusLabels[s]
}
