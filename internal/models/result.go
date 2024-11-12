package models

type ResultRequest struct {
	ChildID        string         `json:"id"`
	ReleaseID      string         `json:"release_id"`
	StageSummaries []StageSummary `json:"stage_summaries"`
	NextStage      string         `json:"next_stage"`
}
