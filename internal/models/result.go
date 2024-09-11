package models

type TimeSummary struct {
	Median  float64
	Minimum float64
	Maximum float64
}

type ResultSummary struct {
	ProxyTimes     TimeSummary
	F1TimesSummary TimeSummary
	F2TimesSummary TimeSummary
	F1ErrRate      float64
	F2ErrRate      float64
}

type ResultRequest struct {
	ID             string        `json:"id"`
	ReleaseID      string        `json:"release_id"`
	ReleaseSummary ResultSummary `json:"release_summary"`
}
