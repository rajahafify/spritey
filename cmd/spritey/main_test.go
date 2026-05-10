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

func decodeCatalogResponse(t *testing.T, data []byte) catalogTestResponse {
	t.Helper()

	var got catalogTestResponse
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to decode JSON %q: %v", string(data), err)
	}
	return got
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
