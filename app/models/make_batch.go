package models

type MakeBatchManifestV1 struct {
	SchemaVersion string                 `json:"schema_version"`
	Jobs          []MakeBatchManifestJob `json:"jobs"`
}

type MakeBatchManifestJob struct {
	ID     string `json:"id"`
	Recipe string `json:"recipe"`
	Out    string `json:"out"`
	Report string `json:"report,omitempty"`
}

type MakeBatchResult struct {
	Command  string               `json:"command"`
	Summary  MakeBatchSummary     `json:"summary"`
	Jobs     []MakeBatchJobResult `json:"jobs"`
	Warnings []string             `json:"warnings"`
}

type MakeBatchSummary struct {
	JobCount     int `json:"job_count"`
	SuccessCount int `json:"success_count"`
	FailedCount  int `json:"failed_count"`
}

type MakeBatchJobResult struct {
	ID       string        `json:"id"`
	Recipe   string        `json:"recipe"`
	Outputs  MakeOutputs   `json:"outputs"`
	Summary  MakeSummary   `json:"summary"`
	Warnings []string      `json:"warnings"`
	Errors   []MakeProblem `json:"errors"`
}
