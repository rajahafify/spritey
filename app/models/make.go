package models

type MakeResult struct {
	Command  string      `json:"command"`
	Outputs  MakeOutputs `json:"outputs"`
	Summary  MakeSummary `json:"summary"`
	Warnings []string    `json:"warnings"`
}

type MakeOutputs struct {
	PNG    MakeFileOutput  `json:"png"`
	Report *MakeFileOutput `json:"report,omitempty"`
}

type MakeFileOutput struct {
	Path string `json:"path"`
}

type MakeSummary struct {
	FrameCount     int        `json:"frame_count"`
	Canvas         MakeCanvas `json:"canvas"`
	AnimationCount int        `json:"animation_count"`
}

type MakeCanvas struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type MakeProblem struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Field   string                 `json:"field,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type MakeReportV1 struct {
	SchemaVersion string `json:"schema_version"`
	Command       string `json:"command"`
	Recipe        struct {
		Path string `json:"path"`
	} `json:"recipe"`
	Assets struct {
		Path string `json:"path"`
	} `json:"assets"`
	Output struct {
		PNG MakeFileOutput `json:"png"`
	} `json:"output"`
	Render struct {
		Canvas struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"canvas"`
		FrameCount   int      `json:"frame_count"`
		AnimationIDs []string `json:"animation_ids"`
	} `json:"render"`
	Layers struct {
		Applied []string `json:"applied"`
	} `json:"layers"`
	Warnings []string `json:"warnings"`
}
