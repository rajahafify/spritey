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

type AssetsValidationResult struct {
	Assets   AssetsValidationTarget  `json:"assets"`
	Pack     Pack                    `json:"pack"`
	Summary  AssetsValidationSummary `json:"summary"`
	Warnings []string                `json:"warnings"`
}

type AssetsValidationTarget struct {
	Path string `json:"path"`
}

type AssetsValidationSummary struct {
	CategoryCount int `json:"category_count"`
	LayerCount    int `json:"layer_count"`
}

type Category struct {
	ID     string  `json:"id"`
	Layers []Layer `json:"layers"`
}

type Layer struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	ZPos            int               `json:"z_pos"`
	BodyTypes       []string          `json:"body_types"`
	BodyTypePaths   map[string]string `json:"-"`
	Animations      []string          `json:"animations"`
	RecolorMaterial string            `json:"recolor_material,omitempty"`
	PathPrefix      string            `json:"path_prefix,omitempty"`
	Credits         []Credit          `json:"credits"`
}

type InspectLayerResult struct {
	Category        string   `json:"category"`
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	ZPos            int      `json:"z_pos"`
	BodyTypes       []string `json:"body_types"`
	Animations      []string `json:"animations"`
	RecolorMaterial string   `json:"recolor_material,omitempty"`
	PathPrefix      string   `json:"path_prefix,omitempty"`
	Credits         []Credit `json:"credits"`
}

type Recipe struct {
	BodyType   string                     `json:"body_type"`
	Selections map[string]RecipeSelection `json:"selections"`
}

type RecipeSelection struct {
	ID             string `json:"id"`
	PaletteVariant string `json:"palette_variant,omitempty"`
}

type RecipeValidationResult struct {
	Path                 string                      `json:"path"`
	BodyTypeRequested    string                      `json:"-"`
	BodyType             string                      `json:"body_type"`
	Selections           []RecipeValidationSelection `json:"selections"`
	Warnings             []string                    `json:"-"`
	RequiredAnimationIDs []string                    `json:"-"`
	RenderInputs         []RecipeRenderInput         `json:"-"`
}

type RecipeValidationSelection struct {
	Category       string `json:"category"`
	ID             string `json:"id"`
	PaletteVariant string `json:"palette_variant,omitempty"`
}

type RecipeRenderInput struct {
	Category         string `json:"-"`
	LayerID          string `json:"-"`
	ResolvedPath     string `json:"-"`
	ResolvedBodyType string `json:"-"`
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
