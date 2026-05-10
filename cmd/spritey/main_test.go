package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type catalogTestResponse struct {
	OK   bool `json:"ok"`
	Pack *struct {
		SchemaVersion string `json:"schema_version"`
		ID            string `json:"id"`
		Name          string `json:"name"`
	} `json:"pack"`
	Categories []struct {
		ID     string `json:"id"`
		Layers []struct {
			ID              string   `json:"id"`
			Name            string   `json:"name"`
			ZPos            int      `json:"z_pos"`
			BodyTypes       []string `json:"body_types"`
			Animations      []string `json:"animations"`
			RecolorMaterial string   `json:"recolor_material,omitempty"`
			PathPrefix      string   `json:"path_prefix,omitempty"`
		} `json:"layers"`
	} `json:"categories"`
	Warnings []string `json:"warnings"`
	Errors   []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Field   string `json:"field,omitempty"`
	} `json:"errors"`
}

type inspectLayerTestResponse struct {
	OK    bool `json:"ok"`
	Layer *struct {
		Category        string   `json:"category"`
		ID              string   `json:"id"`
		Name            string   `json:"name"`
		ZPos            int      `json:"z_pos"`
		BodyTypes       []string `json:"body_types"`
		Animations      []string `json:"animations"`
		RecolorMaterial string   `json:"recolor_material,omitempty"`
		PathPrefix      string   `json:"path_prefix,omitempty"`
		Credits         []struct {
			File     string   `json:"file,omitempty"`
			Authors  []string `json:"authors,omitempty"`
			Licenses []string `json:"licenses,omitempty"`
		} `json:"credits"`
	} `json:"layer"`
	Warnings []string `json:"warnings"`
	Errors   []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Field   string `json:"field,omitempty"`
	} `json:"errors"`
}

type validateRecipeTestResponse struct {
	OK     bool `json:"ok"`
	Recipe *struct {
		Path       string `json:"path"`
		BodyType   string `json:"body_type"`
		Selections []struct {
			Category       string `json:"category"`
			ID             string `json:"id"`
			PaletteVariant string `json:"palette_variant,omitempty"`
		} `json:"selections"`
	} `json:"recipe"`
	Warnings []string `json:"warnings"`
	Errors   []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Field   string `json:"field,omitempty"`
	} `json:"errors"`
}

func TestCatalogJSONSuccess(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets := filepath.Join("..", "..", "testdata", "fixtures", "basic-assets")

	code := run([]string{"catalog", "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}

	got := decodeCatalogResponse(t, stdout.Bytes())
	if !got.OK {
		t.Fatalf("expected ok response: %+v", got)
	}
	if got.Pack == nil || got.Pack.ID != "basic-test-assets" || got.Pack.Name != "Basic Test Assets" {
		t.Fatalf("unexpected pack: %+v", got.Pack)
	}
	if len(got.Errors) != 0 {
		t.Fatalf("expected no errors, got %+v", got.Errors)
	}
	if len(got.Categories) != 2 {
		t.Fatalf("expected 2 categories, got %+v", got.Categories)
	}
	if got.Categories[0].ID != "body" || got.Categories[1].ID != "weapon" {
		t.Fatalf("categories should be sorted by id, got %+v", got.Categories)
	}

	bodyLayer := got.Categories[0].Layers[0]
	if bodyLayer.ID != "body_human" || bodyLayer.Name != "Human Body" || bodyLayer.ZPos != 10 {
		t.Fatalf("unexpected body layer: %+v", bodyLayer)
	}
	assertStrings(t, bodyLayer.BodyTypes, []string{"female", "male"})
	assertStrings(t, bodyLayer.Animations, []string{"walk"})
	if bodyLayer.RecolorMaterial != "skin" {
		t.Fatalf("expected recolor material skin, got %q", bodyLayer.RecolorMaterial)
	}
	if bodyLayer.PathPrefix != "body/human/female/" {
		t.Fatalf("expected first sorted body-type path prefix, got %q", bodyLayer.PathPrefix)
	}
}

func TestCatalogJSONMissingAssets(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"catalog", "--json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}

	got := decodeCatalogResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "MISSING_ASSETS" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestCatalogJSONMissingPack(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets := t.TempDir()

	code := run([]string{"catalog", "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("expected exit code 3, got %d", code)
	}

	got := decodeCatalogResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "MISSING_PACK_JSON" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestCatalogJSONAssetsDirectoryNotFound(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets := filepath.Join(t.TempDir(), "missing-assets")

	code := run([]string{"catalog", "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("expected exit code 3, got %d", code)
	}

	got := decodeCatalogResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "ASSETS_DIRECTORY_NOT_FOUND" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestCatalogJSONInvalidSheetDefinition(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets := t.TempDir()
	writeTestFile(t, filepath.Join(assets, "pack.json"), `{"schema_version":"1","id":"bad-sheet-assets","name":"Bad Sheet Assets"}`)
	writeTestFile(t, filepath.Join(assets, "sheet_definitions", "bad.json"), `{not-json}`)

	code := run([]string{"catalog", "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("expected exit code 3, got %d", code)
	}

	got := decodeCatalogResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "INVALID_SHEET_DEFINITION_JSON" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestInspectLayerJSONSuccess(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets := filepath.Join("..", "..", "testdata", "fixtures", "basic-assets")

	code := run([]string{"inspect", "layer", "body_human", "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}

	got := decodeInspectLayerResponse(t, stdout.Bytes())
	if !got.OK {
		t.Fatalf("expected ok response: %+v", got)
	}
	if got.Layer == nil {
		t.Fatal("expected layer payload")
	}
	if got.Layer.Category != "body" || got.Layer.ID != "body_human" || got.Layer.Name != "Human Body" {
		t.Fatalf("unexpected layer: %+v", got.Layer)
	}
	if got.Layer.ZPos != 10 || got.Layer.RecolorMaterial != "skin" || got.Layer.PathPrefix != "body/human/female/" {
		t.Fatalf("unexpected layer details: %+v", got.Layer)
	}
	assertStrings(t, got.Layer.BodyTypes, []string{"female", "male"})
	assertStrings(t, got.Layer.Animations, []string{"walk"})
	if len(got.Errors) != 0 {
		t.Fatalf("expected no errors, got %+v", got.Errors)
	}
}

func TestInspectLayerJSONMissingTarget(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"inspect", "--json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}

	got := decodeInspectLayerResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "MISSING_INSPECT_TARGET" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestInspectLayerJSONUnsupportedTarget(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"inspect", "pack", "--json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}

	got := decodeInspectLayerResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "UNSUPPORTED_INSPECT_TARGET" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestInspectLayerJSONMissingLayerID(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"inspect", "layer", "--json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}

	got := decodeInspectLayerResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "MISSING_LAYER_ID" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestInspectLayerJSONUnknownLayerID(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets := filepath.Join("..", "..", "testdata", "fixtures", "basic-assets")

	code := run([]string{"inspect", "layer", "missing_layer", "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 5 {
		t.Fatalf("expected exit code 5, got %d", code)
	}

	got := decodeInspectLayerResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "UNKNOWN_LAYER_ID" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestInspectLayerJSONMissingAssets(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"inspect", "layer", "body_human", "--json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}

	got := decodeInspectLayerResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "MISSING_ASSETS" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestInspectLayerJSONInvalidAssetsDirectory(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets := filepath.Join(t.TempDir(), "missing-assets")

	code := run([]string{"inspect", "layer", "body_human", "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("expected exit code 3, got %d", code)
	}

	got := decodeInspectLayerResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "ASSETS_DIRECTORY_NOT_FOUND" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestValidateRecipeJSONSuccess(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets := fixturePath("basic-assets")
	recipe := filepath.Join("..", "..", "testdata", "fixtures", "recipes", "valid-basic.json")

	code := run([]string{"validate", recipe, "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}

	got := decodeValidateRecipeResponse(t, stdout.Bytes())
	if !got.OK {
		t.Fatalf("expected ok response: %+v", got)
	}
	if got.Recipe == nil || got.Recipe.BodyType != "male" {
		t.Fatalf("unexpected recipe: %+v", got.Recipe)
	}
	if len(got.Recipe.Selections) != 2 {
		t.Fatalf("expected 2 selections, got %+v", got.Recipe.Selections)
	}
	if got.Recipe.Selections[0].Category != "body" || got.Recipe.Selections[0].ID != "body_human" {
		t.Fatalf("unexpected first selection: %+v", got.Recipe.Selections[0])
	}
	if len(got.Errors) != 0 {
		t.Fatalf("expected no errors, got %+v", got.Errors)
	}
}

func TestValidateRecipeJSONUsesDefaultBodyType(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets := fixturePath("basic-assets")
	recipe := filepath.Join("..", "..", "testdata", "fixtures", "recipes", "default-body-type.json")

	code := run([]string{"validate", recipe, "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stdout=%q", code, stdout.String())
	}

	got := decodeValidateRecipeResponse(t, stdout.Bytes())
	if got.Recipe == nil || got.Recipe.BodyType != "male" {
		t.Fatalf("expected default body type male, got %+v", got.Recipe)
	}
}

func TestValidateRecipeJSONMissingRecipe(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"validate", "--json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}

	got := decodeValidateRecipeResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "MISSING_RECIPE" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestValidateRecipeJSONMissingAssets(t *testing.T) {
	var stdout, stderr bytes.Buffer
	recipe := filepath.Join("..", "..", "testdata", "fixtures", "recipes", "valid-basic.json")

	code := run([]string{"validate", recipe, "--json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}

	got := decodeValidateRecipeResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "MISSING_ASSETS" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestValidateRecipeJSONMissingRecipeFile(t *testing.T) {
	var stdout, stderr bytes.Buffer
	recipe := filepath.Join(t.TempDir(), "missing.json")

	code := run([]string{"validate", recipe, "--assets", fixturePath("basic-assets"), "--json"}, &stdout, &stderr)
	if code != 4 {
		t.Fatalf("expected exit code 4, got %d", code)
	}

	got := decodeValidateRecipeResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "RECIPE_FILE_NOT_FOUND" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestValidateRecipeJSONInvalidRecipeJSON(t *testing.T) {
	var stdout, stderr bytes.Buffer
	recipe := filepath.Join("..", "..", "testdata", "fixtures", "recipes", "invalid-json.json")

	code := run([]string{"validate", recipe, "--assets", fixturePath("basic-assets"), "--json"}, &stdout, &stderr)
	if code != 4 {
		t.Fatalf("expected exit code 4, got %d", code)
	}

	got := decodeValidateRecipeResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "INVALID_RECIPE_JSON" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestValidateRecipeJSONMissingSelections(t *testing.T) {
	var stdout, stderr bytes.Buffer
	recipe := filepath.Join("..", "..", "testdata", "fixtures", "recipes", "missing-selections.json")

	code := run([]string{"validate", recipe, "--assets", fixturePath("basic-assets"), "--json"}, &stdout, &stderr)
	if code != 4 {
		t.Fatalf("expected exit code 4, got %d", code)
	}

	got := decodeValidateRecipeResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "MISSING_SELECTIONS" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestValidateRecipeJSONUnknownLayer(t *testing.T) {
	var stdout, stderr bytes.Buffer
	recipe := filepath.Join("..", "..", "testdata", "fixtures", "recipes", "unknown-layer.json")

	code := run([]string{"validate", recipe, "--assets", fixturePath("basic-assets"), "--json"}, &stdout, &stderr)
	if code != 5 {
		t.Fatalf("expected exit code 5, got %d", code)
	}

	got := decodeValidateRecipeResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "UNKNOWN_LAYER_ID" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestValidateRecipeJSONUnsupportedBodyType(t *testing.T) {
	var stdout, stderr bytes.Buffer
	recipe := filepath.Join("..", "..", "testdata", "fixtures", "recipes", "unsupported-body-type.json")

	code := run([]string{"validate", recipe, "--assets", fixturePath("basic-assets"), "--json"}, &stdout, &stderr)
	if code != 5 {
		t.Fatalf("expected exit code 5, got %d", code)
	}

	got := decodeValidateRecipeResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "UNSUPPORTED_BODY_TYPE" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func decodeCatalogResponse(t *testing.T, data []byte) catalogTestResponse {
	t.Helper()

	var got catalogTestResponse
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to decode JSON %q: %v", string(data), err)
	}
	return got
}

func decodeInspectLayerResponse(t *testing.T, data []byte) inspectLayerTestResponse {
	t.Helper()

	var got inspectLayerTestResponse
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to decode JSON %q: %v", string(data), err)
	}
	return got
}

func decodeValidateRecipeResponse(t *testing.T, data []byte) validateRecipeTestResponse {
	t.Helper()

	var got validateRecipeTestResponse
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to decode JSON %q: %v", string(data), err)
	}
	return got
}

func fixturePath(name string) string {
	return filepath.Join("..", "..", "testdata", "fixtures", name)
}

func assertStrings(t *testing.T, got []string, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, got)
		}
	}
}

func writeTestFile(t *testing.T, path string, data string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}
