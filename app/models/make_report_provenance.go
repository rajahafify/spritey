package models

type MakeReportV1Provenance struct {
	SchemaVersion string `json:"schema_version"`
	Command       string `json:"command"`
	Pack          struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"pack"`
	Recipe struct {
		Path              string `json:"path"`
		BodyTypeEffective string `json:"body_type_effective"`
		BodyTypeRequested string `json:"body_type_requested"`
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
		Applied  []string                  `json:"applied"`
		Composed []MakeReportComposedLayer `json:"composed"`
	} `json:"layers"`
	Warnings []string `json:"warnings"`
}

type MakeReportComposedLayer struct {
	Category         string   `json:"category"`
	ID               string   `json:"id"`
	ZPos             int      `json:"z_pos"`
	ResolvedBodyType string   `json:"resolved_body_type"`
	ResolvedPath     string   `json:"resolved_path"`
	PaletteVariant   string   `json:"palette_variant"`
	Credits          []Credit `json:"credits"`
}
