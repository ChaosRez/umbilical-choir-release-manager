package models

type TimeSummary struct {
	Median  float64 `json:"median"`
	Minimum float64 `json:"minimum"`
	Maximum float64 `json:"maximum"`
}

type ResultSummary struct {
	StageName      string      `json:"stage_name"`
	ProxyTimes     TimeSummary `json:"proxy_times"`
	F1TimesSummary TimeSummary `json:"f1_times_summary"`
	F2TimesSummary TimeSummary `json:"f2_times_summary"`
	F1ErrRate      float64     `json:"f1_err_rate"`
	F2ErrRate      float64     `json:"f2_err_rate"`
}

type ResultRequest struct {
	ChildID          string          `json:"id"`
	ReleaseID        string          `json:"release_id"`
	ReleaseSummaries []ResultSummary `json:"release_summaries"`
}
