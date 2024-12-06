package models

import "fmt"

type StageSummary struct {
	StageName      string      `json:"stage_name"`
	ProxyTimes     TimeSummary `json:"proxy_times"`
	F1TimesSummary TimeSummary `json:"f1_times_summary"`
	F2TimesSummary TimeSummary `json:"f2_times_summary"`
	F1ErrRate      float64     `json:"f1_err_rate"`
	F2ErrRate      float64     `json:"f2_err_rate"`
	Status         StageStatus `json:"status"`
}

type TimeSummary struct {
	Median  float64 `json:"median"`
	Minimum float64 `json:"minimum"`
	Maximum float64 `json:"maximum"`
}

type StageStatus int

type StageStatusKey struct {
	ReleaseID string
	StageName string
	ChildID   string
}

const ( // NOTE, for any change, update stageStatusLabels (+ the readme, and agent)
	Pending        StageStatus = iota // The first stage status will be initialized as InProgress and never will be Pending
	InProgress                        // the child is notified
	SuccessWaiting                    // only WaitForSignal stage type. The child received enough calls and was successful
	ShouldEnd                         // only WaitForSignal stage type. The child poll for it on /end_stage to finish a stage (set by the parent)
	Completed                         // received the stage result
	Failure                           // received the stage result as Failure
	Error                             // received the stage result as Error
)

var stageStatusLabels = []string{
	"Pending",
	"InProgress",
	"SuccessWaiting",
	"ShouldEnd",
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
