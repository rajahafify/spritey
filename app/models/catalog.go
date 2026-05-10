package models

type Pack struct {
	SchemaVersion string       `json:"schema_version"`
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Defaults      PackDefaults `json:"defaults,omitempty"`
}

type PackDefaults struct {
	BodyType                string   `json:"body_type,omitempty"`
	Animations              []string `json:"animations,omitempty"`
	CanvasWidth             int      `json:"canvas_width,omitempty"`
	OutputFormat            string   `json:"output_format,omitempty"`
	MissingBodyTypeFallback string   `json:"missing_body_type_fallback,omitempty"`
	PaletteSourceFallbacks  []string `json:"palette_source_fallbacks,omitempty"`
}

type Catalog struct {
	Pack       Pack       `json:"pack"`
	Categories []Category `json:"categories"`
	Warnings   []string   `json:"warnings"`
}

type Category struct {
	ID     string  `json:"id"`
	Layers []Layer `json:"layers"`
}

type Layer struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	ZPos            int      `json:"z_pos"`
	BodyTypes       []string `json:"body_types"`
	Animations      []string `json:"animations"`
	RecolorMaterial string   `json:"recolor_material,omitempty"`
	PathPrefix      string   `json:"path_prefix,omitempty"`
	Credits         []Credit `json:"credits"`
}

type Credit struct {
	File     string   `json:"file,omitempty"`
	Notes    string   `json:"notes,omitempty"`
	Authors  []string `json:"authors,omitempty"`
	Licenses []string `json:"licenses,omitempty"`
	URLs     []string `json:"urls,omitempty"`
}

type Problem struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func (p *Problem) Error() string {
	if p == nil {
		return ""
	}
	return p.Message
}
