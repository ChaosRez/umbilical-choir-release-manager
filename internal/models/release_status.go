package models

import "fmt"

type ReleaseStatus int

const (
	No    ReleaseStatus = iota // the child should not get the release instruction
	Todo                       // Marked to get the release
	Doing                      // the child is notified of the release
	Done                       // the child done all stages
	Failed
)

var releaseStatusLabels = []string{
	"No",
	"Todo",
	"Doing",
	"Done",
	"Failed",
}

// String returns the string representation of the ReleaseStatus (you can print as %s)
func (rs ReleaseStatus) String() string {
	if rs < 0 || int(rs) >= len(releaseStatusLabels) {
		return fmt.Sprintf("ReleaseStatus(%d)", rs)
	}
	return releaseStatusLabels[rs]
}
