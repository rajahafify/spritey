package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rajahafify/spritey/app/services"
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

type assetsValidationTestResponse struct {
	OK     bool `json:"ok"`
	Assets *struct {
		Path string `json:"path"`
	} `json:"assets"`
	Pack *struct {
		SchemaVersion string `json:"schema_version"`
		ID            string `json:"id"`
		Name          string `json:"name"`
	} `json:"pack"`
	Summary *struct {
		CategoryCount int `json:"category_count"`
		LayerCount    int `json:"layer_count"`
	} `json:"summary"`
	Warnings []string `json:"warnings"`
	Errors   []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Field   string `json:"field,omitempty"`
	} `json:"errors"`
}

type makeTestResponse struct {
	OK      bool   `json:"ok"`
	Command string `json:"command"`
	Outputs struct {
		PNG *struct {
			Path string `json:"path"`
		} `json:"png"`
		Report *struct {
			Path string `json:"path"`
		} `json:"report,omitempty"`
	} `json:"outputs"`
	Summary  map[string]interface{} `json:"summary"`
	Warnings []string               `json:"warnings"`
	Errors   []struct {
		Code    string                 `json:"code"`
		Message string                 `json:"message"`
		Field   string                 `json:"field,omitempty"`
		Details map[string]interface{} `json:"details,omitempty"`
	} `json:"errors"`
}

type makeBatchTestResponse struct {
	OK      bool   `json:"ok"`
	Command string `json:"command"`
	Summary struct {
		JobCount     int `json:"job_count"`
		SuccessCount int `json:"success_count"`
		FailedCount  int `json:"failed_count"`
	} `json:"summary"`
	Jobs []struct {
		ID      string `json:"id"`
		Recipe  string `json:"recipe"`
		Outputs struct {
			PNG *struct {
				Path string `json:"path"`
			} `json:"png"`
			Report *struct {
				Path string `json:"path"`
			} `json:"report,omitempty"`
		} `json:"outputs"`
		Summary  map[string]interface{} `json:"summary"`
		Warnings []string               `json:"warnings"`
		Errors   []struct {
			Code    string                 `json:"code"`
			Message string                 `json:"message"`
			Field   string                 `json:"field,omitempty"`
			Details map[string]interface{} `json:"details,omitempty"`
		} `json:"errors"`
	} `json:"jobs"`
	Warnings []string `json:"warnings"`
	Errors   []struct {
		Code    string                 `json:"code"`
		Message string                 `json:"message"`
		Field   string                 `json:"field,omitempty"`
		Details map[string]interface{} `json:"details,omitempty"`
	} `json:"errors"`
}

type downloadAssetsTestResponse struct {
	OK      bool   `json:"ok"`
	Command string `json:"command"`
	Assets  *struct {
		Path    string `json:"path"`
		Source  string `json:"source,omitempty"`
		Version string `json:"version,omitempty"`
	} `json:"assets"`
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

func TestDownloadLPCAssetsJSONSuccess(t *testing.T) {
	var stdout, stderr bytes.Buffer
	tempCache := t.TempDir()
	zipPayload := createLPCArchiveZipPayload(t, false)

	originalFetcher := services.DefaultLPCArchiveFetcher
	originalCacheResolver := services.DefaultUserCacheDirResolver
	services.DefaultLPCArchiveFetcher = func(url string, onProgress func(int64, int64)) ([]byte, error) {
		return zipPayload, nil
	}
	services.DefaultUserCacheDirResolver = func() (string, error) {
		return tempCache, nil
	}
	t.Cleanup(func() {
		services.DefaultLPCArchiveFetcher = originalFetcher
		services.DefaultUserCacheDirResolver = originalCacheResolver
	})

	code := run([]string{"--download-lpc-assets", "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	got := decodeDownloadAssetsResponse(t, stdout.Bytes())
	if !got.OK || got.Command != "download-lpc-assets" {
		t.Fatalf("unexpected response: %+v", got)
	}
	if got.Assets == nil {
		t.Fatalf("expected assets payload, got %+v", got)
	}
	expectedPath := filepath.Join(tempCache, "spritey", "assets", "lpc")
	if got.Assets.Path != expectedPath {
		t.Fatalf("expected path %q, got %q", expectedPath, got.Assets.Path)
	}
	if len(got.Errors) != 0 {
		t.Fatalf("expected no errors, got %+v", got.Errors)
	}
}

func TestDownloadLPCAssetsUnknownArgumentIsCLIUsage(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"--download-lpc-assets", "--json", "--wat"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}

	got := decodeDownloadAssetsResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "UNKNOWN_ARGUMENT" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestDownloadLPCAssetsDownloadFailureMapsToExit1(t *testing.T) {
	var stdout, stderr bytes.Buffer
	tempCache := t.TempDir()

	originalFetcher := services.DefaultLPCArchiveFetcher
	originalCacheResolver := services.DefaultUserCacheDirResolver
	services.DefaultLPCArchiveFetcher = func(url string, onProgress func(int64, int64)) ([]byte, error) {
		return nil, fmt.Errorf("network unavailable")
	}
	services.DefaultUserCacheDirResolver = func() (string, error) {
		return tempCache, nil
	}
	t.Cleanup(func() {
		services.DefaultLPCArchiveFetcher = originalFetcher
		services.DefaultUserCacheDirResolver = originalCacheResolver
	})

	code := run([]string{"--download-lpc-assets", "--json"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	got := decodeDownloadAssetsResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "DOWNLOAD_FAILED" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestDownloadLPCAssetsInvalidDownloadedAssetsMapsToExit3(t *testing.T) {
	var stdout, stderr bytes.Buffer
	tempCache := t.TempDir()
	zipPayload := createLPCArchiveZipPayload(t, true)

	originalFetcher := services.DefaultLPCArchiveFetcher
	originalCacheResolver := services.DefaultUserCacheDirResolver
	services.DefaultLPCArchiveFetcher = func(url string, onProgress func(int64, int64)) ([]byte, error) {
		return zipPayload, nil
	}
	services.DefaultUserCacheDirResolver = func() (string, error) {
		return tempCache, nil
	}
	t.Cleanup(func() {
		services.DefaultLPCArchiveFetcher = originalFetcher
		services.DefaultUserCacheDirResolver = originalCacheResolver
	})

	code := run([]string{"--download-lpc-assets", "--json"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("expected exit code 3, got %d", code)
	}
	got := decodeDownloadAssetsResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "INVALID_PACK_JSON" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestDownloadLPCAssetsShowsProgressInTextMode(t *testing.T) {
	var stdout, stderr bytes.Buffer
	tempCache := t.TempDir()
	zipPayload := createLPCArchiveZipPayload(t, false)

	originalFetcher := services.DefaultLPCArchiveFetcher
	originalCacheResolver := services.DefaultUserCacheDirResolver
	services.DefaultLPCArchiveFetcher = func(url string, onProgress func(int64, int64)) ([]byte, error) {
		if onProgress != nil {
			onProgress(0, 100)
			onProgress(25, 100)
			onProgress(75, 100)
			onProgress(100, 100)
		}
		return zipPayload, nil
	}
	services.DefaultUserCacheDirResolver = func() (string, error) {
		return tempCache, nil
	}
	t.Cleanup(func() {
		services.DefaultLPCArchiveFetcher = originalFetcher
		services.DefaultUserCacheDirResolver = originalCacheResolver
	})

	code := run([]string{"--download-lpc-assets"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if !strings.Contains(stdout.String(), "ok: ") {
		t.Fatalf("expected success text output, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Downloading LPC assets") {
		t.Fatalf("expected progress output on stderr, got %q", stderr.String())
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
	assets := t.TempDir()
	writeTestFile(t, filepath.Join(assets, "pack.json"), `{"schema_version":"1","id":"unsupported-assets","name":"Unsupported Assets","defaults":{"body_type":"male","animations":["idle"],"canvas_width":8}}`)
	writeTestFile(t, filepath.Join(assets, "sheet_definitions", "body", "body_human.json"), `{"name":"Human Body","type_name":"body","layer_1":{"zPos":10,"male":"body/human/male/"},"animations":["idle"]}`)
	writeTestFile(t, filepath.Join(assets, "palette_definitions", ".gitkeep"), "")
	writeTestFile(t, filepath.Join(assets, "spritesheets", ".gitkeep"), "")
	makeDeterministicLayerPNG(t, filepath.Join(assets, "spritesheets", "body", "human", "male", "idle.png"), color.RGBA{R: 100, G: 120, B: 220, A: 255})

	recipe := filepath.Join(t.TempDir(), "unsupported.json")
	writeTestFile(t, recipe, `{"body_type":"child","selections":{"body":{"id":"body_human"}}}`)

	code := run([]string{"validate", recipe, "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 5 {
		t.Fatalf("expected exit code 5, got %d", code)
	}

	got := decodeValidateRecipeResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "UNSUPPORTED_BODY_TYPE" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestValidateRecipeJSONMissingSpriteFrame(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, recipe := writeReadinessCLIFixture(t, readinessCLIFixtureOptions{
		recipeBodyType:      "male",
		missingFallback:     "male",
		requiredAnimations:  []string{"idle", "walk"},
		includeBodyIdle:     true,
		includeBodyWalk:     true,
		includeWeaponIdle:   true,
		includeWeaponWalk:   false,
		weaponHasFemalePath: true,
	})

	code := run([]string{"validate", recipe, "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 5 {
		t.Fatalf("expected exit code 5, got %d", code)
	}
	got := decodeValidateRecipeResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "MISSING_SPRITE_FRAME" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestValidateRecipeJSONFallbackPathWarning(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, recipe := writeReadinessCLIFixture(t, readinessCLIFixtureOptions{
		recipeBodyType:        "female",
		missingFallback:       "male",
		requiredAnimations:    []string{"idle"},
		includeBodyIdle:       true,
		includeBodyFemaleIdle: true,
		includeWeaponIdle:     true,
		weaponHasFemalePath:   false,
	})

	code := run([]string{"validate", recipe, "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	got := decodeValidateRecipeResponse(t, stdout.Bytes())
	if !got.OK {
		t.Fatalf("expected success response: %+v", got)
	}
	if len(got.Warnings) != 1 {
		t.Fatalf("expected one warning, got %+v", got.Warnings)
	}
}

func TestAssetsValidateJSONSuccess(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets := fixturePath("basic-assets")

	code := run([]string{"assets", "validate", "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}

	got := decodeAssetsValidationResponse(t, stdout.Bytes())
	if !got.OK {
		t.Fatalf("expected ok response: %+v", got)
	}
	if got.Assets == nil || got.Assets.Path != assets {
		t.Fatalf("unexpected assets payload: %+v", got.Assets)
	}
	if got.Pack == nil || got.Pack.ID != "basic-test-assets" || got.Pack.Name != "Basic Test Assets" {
		t.Fatalf("unexpected pack: %+v", got.Pack)
	}
	if got.Summary == nil || got.Summary.CategoryCount != 2 || got.Summary.LayerCount != 2 {
		t.Fatalf("unexpected summary: %+v", got.Summary)
	}
	if len(got.Warnings) != 0 || len(got.Errors) != 0 {
		t.Fatalf("expected no warnings or errors, got warnings=%+v errors=%+v", got.Warnings, got.Errors)
	}
}

func TestAssetsValidateJSONMissingSubcommand(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"assets", "--json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}

	got := decodeAssetsValidationResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "MISSING_ASSETS_SUBCOMMAND" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestAssetsValidateJSONUnsupportedSubcommand(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"assets", "inspect", "--json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}

	got := decodeAssetsValidationResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "UNSUPPORTED_ASSETS_SUBCOMMAND" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestAssetsValidateJSONMissingAssets(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"assets", "validate", "--json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}

	got := decodeAssetsValidationResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "MISSING_ASSETS" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestAssetsValidateJSONAssetsDirectoryNotFound(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets := filepath.Join(t.TempDir(), "missing-assets")

	code := run([]string{"assets", "validate", "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("expected exit code 3, got %d", code)
	}

	got := decodeAssetsValidationResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "ASSETS_DIRECTORY_NOT_FOUND" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestAssetsValidateJSONMissingPack(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets := t.TempDir()

	code := run([]string{"assets", "validate", "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("expected exit code 3, got %d", code)
	}

	got := decodeAssetsValidationResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "MISSING_PACK_JSON" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestAssetsValidateJSONInvalidPack(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets := t.TempDir()
	writeTestFile(t, filepath.Join(assets, "pack.json"), `{not-json}`)

	code := run([]string{"assets", "validate", "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("expected exit code 3, got %d", code)
	}

	got := decodeAssetsValidationResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "INVALID_PACK_JSON" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestAssetsValidateJSONMissingSheetDefinitions(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets := t.TempDir()
	writeTestFile(t, filepath.Join(assets, "pack.json"), `{"schema_version":"1","id":"missing-sheets","name":"Missing Sheets"}`)

	code := run([]string{"assets", "validate", "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("expected exit code 3, got %d", code)
	}

	got := decodeAssetsValidationResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "MISSING_SHEET_DEFINITIONS" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestAssetsValidateJSONInvalidSheetDefinition(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets := t.TempDir()
	writeTestFile(t, filepath.Join(assets, "pack.json"), `{"schema_version":"1","id":"bad-sheet-assets","name":"Bad Sheet Assets"}`)
	writeTestFile(t, filepath.Join(assets, "sheet_definitions", "bad.json"), `{not-json}`)

	code := run([]string{"assets", "validate", "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("expected exit code 3, got %d", code)
	}

	got := decodeAssetsValidationResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "INVALID_SHEET_DEFINITION_JSON" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestAssetsValidateJSONMissingSpritesheets(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets := t.TempDir()
	writeMinimalAssetsPack(t, assets)
	if err := os.Remove(filepath.Join(assets, "spritesheets", ".gitkeep")); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(assets, "spritesheets")); err != nil {
		t.Fatal(err)
	}

	code := run([]string{"assets", "validate", "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("expected exit code 3, got %d", code)
	}

	got := decodeAssetsValidationResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "MISSING_SPRITESHEETS" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestAssetsValidateJSONMissingPaletteDefinitions(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets := t.TempDir()
	writeMinimalAssetsPack(t, assets)
	if err := os.Remove(filepath.Join(assets, "palette_definitions", ".gitkeep")); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(assets, "palette_definitions")); err != nil {
		t.Fatal(err)
	}

	code := run([]string{"assets", "validate", "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("expected exit code 3, got %d", code)
	}

	got := decodeAssetsValidationResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "MISSING_PALETTE_DEFINITIONS" {
		t.Fatalf("unexpected error response: %+v", got)
	}
}

func TestMakeJSONSuccessWithReport(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, recipe := writeDeterministicMakeFixture(t)
	out := filepath.Join(t.TempDir(), "sprite.png")
	report := filepath.Join(t.TempDir(), "sprite.report.json")

	code := run([]string{"make", recipe, "--assets", assets, "--out", out, "--report", report, "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	got := decodeMakeResponse(t, stdout.Bytes())
	if !got.OK || got.Command != "make" {
		t.Fatalf("unexpected response: %+v", got)
	}
	if got.Outputs.PNG == nil || got.Outputs.PNG.Path != out {
		t.Fatalf("unexpected png output: %+v", got.Outputs)
	}
	if got.Outputs.Report == nil || got.Outputs.Report.Path != report {
		t.Fatalf("unexpected report output: %+v", got.Outputs)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("expected png output to exist: %v", err)
	}
	if _, err := os.Stat(report); err != nil {
		t.Fatalf("expected report output to exist: %v", err)
	}
}

func TestMakeJSONReportProvenanceAndEnvelopeUnchanged(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, recipe := writeDeterministicMakeFixture(t)
	out := filepath.Join(t.TempDir(), "sprite.png")
	report := filepath.Join(t.TempDir(), "sprite.report.json")

	code := run([]string{"make", recipe, "--assets", assets, "--out", out, "--report", report, "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}

	var envelope map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &envelope); err != nil {
		t.Fatal(err)
	}
	assertHasKeys(t, envelope, "ok", "command", "outputs", "summary", "warnings", "errors")
	if len(envelope) != 6 {
		t.Fatalf("expected unchanged make JSON top-level envelope keys, got %+v", envelope)
	}
	summary, ok := envelope["summary"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected summary object in envelope, got %+v", envelope["summary"])
	}
	canvas, ok := summary["canvas"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected summary.canvas object, got %+v", summary["canvas"])
	}
	if int(canvas["width"].(float64)) != 832 || int(canvas["height"].(float64)) != 8 {
		t.Fatalf("expected LPC summary canvas 832x8, got %+v", canvas)
	}
	if int(summary["frame_count"].(float64)) != 1 || int(summary["animation_count"].(float64)) != 1 {
		t.Fatalf("expected emitted-row summary counts, got %+v", summary)
	}

	reportData, err := os.ReadFile(report)
	if err != nil {
		t.Fatal(err)
	}
	var reportJSON map[string]interface{}
	if err := json.Unmarshal(reportData, &reportJSON); err != nil {
		t.Fatal(err)
	}
	pack, ok := reportJSON["pack"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected pack metadata in report: %+v", reportJSON)
	}
	if pack["id"] != "make-assets" || pack["name"] != "Make Assets" {
		t.Fatalf("unexpected report pack metadata: %+v", pack)
	}

	recipeMeta := reportJSON["recipe"].(map[string]interface{})
	if recipeMeta["path"] != recipe || recipeMeta["body_type_requested"] != "male" || recipeMeta["body_type_effective"] != "male" {
		t.Fatalf("unexpected report recipe provenance: %+v", recipeMeta)
	}
	artifacts, ok := reportJSON["artifacts"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected artifacts metadata in report: %+v", reportJSON)
	}
	outputPNG, ok := artifacts["output_png"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected artifacts.output_png metadata in report: %+v", artifacts)
	}
	if outputPNG["sha256"] != pngSHA256(t, out) {
		t.Fatalf("unexpected report artifact sha256: %+v", outputPNG)
	}
	info, err := os.Stat(out)
	if err != nil {
		t.Fatal(err)
	}
	if int64(outputPNG["bytes"].(float64)) != info.Size() {
		t.Fatalf("unexpected report artifact bytes: got=%v want=%d", outputPNG["bytes"], info.Size())
	}
	layers := reportJSON["layers"].(map[string]interface{})
	composed, ok := layers["composed"].([]interface{})
	if !ok || len(composed) != 2 {
		t.Fatalf("expected composed layer provenance in report, got %+v", layers["composed"])
	}
}

func TestMakeReportAnimationOrderMatchesStripRowOrder(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, recipe := writeDeterministicMakeFixture(t)
	out := filepath.Join(t.TempDir(), "sprite.png")
	report := filepath.Join(t.TempDir(), "sprite.report.json")

	code := run([]string{"make", recipe, "--assets", assets, "--out", out, "--report", report, "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}

	reportData, err := os.ReadFile(report)
	if err != nil {
		t.Fatal(err)
	}
	var reportJSON map[string]interface{}
	if err := json.Unmarshal(reportData, &reportJSON); err != nil {
		t.Fatal(err)
	}

	render, ok := reportJSON["render"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected render object in report, got %+v", reportJSON["render"])
	}
	animationIDs, ok := render["animation_ids"].([]interface{})
	if !ok {
		t.Fatalf("expected render.animation_ids array, got %+v", render["animation_ids"])
	}
	if len(animationIDs) != 1 || animationIDs[0] != "walk" {
		t.Fatalf("unexpected animation id order in report: %+v", animationIDs)
	}

	file, err := os.Open(out)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		t.Fatal(err)
	}

	top := color.RGBAModel.Convert(img.At(2, 2)).(color.RGBA)
	if top != (color.RGBA{R: 214, G: 50, B: 90, A: 255}) {
		t.Fatalf("expected only emitted row to match walk frame, got %+v", top)
	}
}

func TestMakeJSONSuccessWithoutReport(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, recipe := writeDeterministicMakeFixture(t)
	out := filepath.Join(t.TempDir(), "sprite.png")

	code := run([]string{"make", recipe, "--assets", assets, "--out", out, "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	got := decodeMakeResponse(t, stdout.Bytes())
	if !got.OK {
		t.Fatalf("unexpected response: %+v", got)
	}
	if got.Outputs.Report != nil {
		t.Fatalf("expected report output to be omitted, got %+v", got.Outputs.Report)
	}
}

func TestMakeTextSuccessWithReport(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, recipe := writeDeterministicMakeFixture(t)
	out := filepath.Join(t.TempDir(), "sprite.png")
	report := filepath.Join(t.TempDir(), "sprite.report.json")

	code := run([]string{"make", recipe, "--assets", assets, "--out", out, "--report", report}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}

	want := fmt.Sprintf(
		"ok: make\npng: %s\nreport: %s\nframe_count: 1\ncanvas: 832x8\nanimation_count: 1\n",
		out,
		report,
	)
	if stdout.String() != want {
		t.Fatalf("unexpected stdout\nwant:\n%s\ngot:\n%s", want, stdout.String())
	}
}

func TestMakeTextSuccessWithoutReport(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, recipe := writeDeterministicMakeFixture(t)
	out := filepath.Join(t.TempDir(), "sprite.png")

	code := run([]string{"make", recipe, "--assets", assets, "--out", out}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}

	want := fmt.Sprintf(
		"ok: make\npng: %s\nframe_count: 1\ncanvas: 832x8\nanimation_count: 1\n",
		out,
	)
	if stdout.String() != want {
		t.Fatalf("unexpected stdout\nwant:\n%s\ngot:\n%s", want, stdout.String())
	}
}

func TestMakeMissingRecipe(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"make", "--assets", "x", "--out", "y", "--json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	got := decodeMakeResponse(t, stdout.Bytes())
	if got.OK || got.Errors[0].Code != "MISSING_RECIPE" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestMakeMissingAssets(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"make", "recipe.json", "--out", "y", "--json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	got := decodeMakeResponse(t, stdout.Bytes())
	if got.OK || got.Errors[0].Code != "MISSING_ASSETS" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestMakeMissingOut(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"make", "recipe.json", "--assets", "x", "--json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	got := decodeMakeResponse(t, stdout.Bytes())
	if got.OK || got.Errors[0].Code != "MISSING_OUT" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestMakeTextMissingOut(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"make", "recipe.json", "--assets", "x"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected empty stdout, got %q", stdout.String())
	}
	if stderr.String() != "MISSING_OUT: --out is required\n" {
		t.Fatalf("unexpected stderr: %q", stderr.String())
	}
}

func TestMakeNonexistentRecipeAndAssets(t *testing.T) {
	var stdout, stderr bytes.Buffer
	out := filepath.Join(t.TempDir(), "sprite.png")
	code := run([]string{"make", filepath.Join(t.TempDir(), "missing.json"), "--assets", filepath.Join(t.TempDir(), "missing-assets"), "--out", out, "--json"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("expected exit code 3, got %d", code)
	}
	got := decodeMakeResponse(t, stdout.Bytes())
	if got.OK || got.Errors[0].Code != "ASSETS_DIRECTORY_NOT_FOUND" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestMakeJSONEnvelopeStability(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, recipe := writeDeterministicMakeFixture(t)
	out := filepath.Join(t.TempDir(), "sprite.png")
	code := run([]string{"make", recipe, "--assets", assets, "--out", out, "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &parsed); err != nil {
		t.Fatal(err)
	}
	assertHasKeys(t, parsed, "ok", "command", "outputs", "summary", "warnings", "errors")
}

func TestMakeJSONErrorEnvelopeStability(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"make", "--json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &parsed); err != nil {
		t.Fatal(err)
	}
	assertHasKeys(t, parsed, "ok", "command", "outputs", "summary", "warnings", "errors")
	errorsField, ok := parsed["errors"].([]interface{})
	if !ok || len(errorsField) == 0 {
		t.Fatalf("expected non-empty errors: %+v", parsed)
	}
	errObj, ok := errorsField[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected error object: %+v", errorsField[0])
	}
	assertHasKeys(t, errObj, "code", "message")
}

func TestMakeRenderFailureIsExit6(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, recipe := writeDeterministicMakeFixture(t)
	outDir := filepath.Join(t.TempDir(), "out-dir")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatal(err)
	}

	code := run([]string{"make", recipe, "--assets", assets, "--out", outDir, "--json"}, &stdout, &stderr)
	if code != 6 {
		t.Fatalf("expected exit code 6, got %d", code)
	}
	got := decodeMakeResponse(t, stdout.Bytes())
	if got.OK || got.Errors[0].Code != "RENDER_FAILED" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestMakeMissingSpriteFrameStillSucceedsWithAvailableRows(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, recipe := writeReadinessCLIFixture(t, readinessCLIFixtureOptions{
		recipeBodyType:      "male",
		missingFallback:     "male",
		requiredAnimations:  []string{"idle", "walk"},
		includeBodyIdle:     true,
		includeBodyWalk:     true,
		includeWeaponIdle:   true,
		includeWeaponWalk:   false,
		weaponHasFemalePath: true,
	})
	out := filepath.Join(t.TempDir(), "sprite.png")
	report := filepath.Join(t.TempDir(), "sprite.report.json")

	code := run([]string{"make", recipe, "--assets", assets, "--out", out, "--report", report, "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	got := decodeMakeResponse(t, stdout.Bytes())
	if !got.OK {
		t.Fatalf("expected success response: %+v", got)
	}
	if int(got.Summary["frame_count"].(float64)) != 1 || int(got.Summary["animation_count"].(float64)) != 1 {
		t.Fatalf("expected one emitted row summary, got %+v", got.Summary)
	}
	if _, err := os.Stat(out); err != nil {
		t.Fatalf("expected png output to exist, stat err=%v", err)
	}
}

func TestMakeSlashUsesMappedAttackPathWhenDirectFramesMissing(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, recipe := writeSlashParityCLIFixture(t, true)
	out := filepath.Join(t.TempDir(), "sprite.png")

	code := run([]string{"make", recipe, "--assets", assets, "--out", out, "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	got := decodeMakeResponse(t, stdout.Bytes())
	if !got.OK {
		t.Fatalf("expected success response, got %+v", got)
	}
	if got.Outputs.PNG == nil || got.Outputs.PNG.Path != out {
		t.Fatalf("unexpected png output: %+v", got.Outputs.PNG)
	}
}

func TestMakeSlashMissingMappedPathStillSucceedsWithBodyRow(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, recipe := writeSlashParityCLIFixture(t, false)
	out := filepath.Join(t.TempDir(), "sprite.png")

	code := run([]string{"make", recipe, "--assets", assets, "--out", out, "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	got := decodeMakeResponse(t, stdout.Bytes())
	if !got.OK {
		t.Fatalf("expected success response: %+v", got)
	}
	if int(got.Summary["frame_count"].(float64)) != 1 || int(got.Summary["animation_count"].(float64)) != 1 {
		t.Fatalf("expected one emitted row summary, got %+v", got.Summary)
	}
}

func TestMakeFallbackWarningInJSONAndReport(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, recipe := writeReadinessCLIFixture(t, readinessCLIFixtureOptions{
		recipeBodyType:        "female",
		missingFallback:       "male",
		requiredAnimations:    []string{"idle"},
		includeBodyIdle:       true,
		includeBodyFemaleIdle: true,
		includeWeaponIdle:     true,
		weaponHasFemalePath:   false,
	})
	out := filepath.Join(t.TempDir(), "sprite.png")
	report := filepath.Join(t.TempDir(), "sprite.report.json")

	code := run([]string{"make", recipe, "--assets", assets, "--out", out, "--report", report, "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	got := decodeMakeResponse(t, stdout.Bytes())
	if !got.OK {
		t.Fatalf("expected success: %+v", got)
	}
	if len(got.Warnings) != 1 {
		t.Fatalf("expected warning in make response, got %+v", got.Warnings)
	}

	reportData, err := os.ReadFile(report)
	if err != nil {
		t.Fatal(err)
	}
	var reportJSON map[string]interface{}
	if err := json.Unmarshal(reportData, &reportJSON); err != nil {
		t.Fatal(err)
	}
	warnings, ok := reportJSON["warnings"].([]interface{})
	if !ok || len(warnings) != 1 {
		t.Fatalf("expected warning in report, got %+v", reportJSON["warnings"])
	}
}

func TestMakeBatchJSONSuccessEnvelopeAndOrdering(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, manifest := writeSuccessfulBatchManifestFixture(t)

	code := run([]string{"make", "batch", manifest, "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &parsed); err != nil {
		t.Fatal(err)
	}
	assertHasKeys(t, parsed, "ok", "command", "summary", "jobs", "warnings", "errors")
	if len(parsed) != 6 {
		t.Fatalf("expected stable top-level key set, got %+v", parsed)
	}

	got := decodeMakeBatchResponse(t, stdout.Bytes())
	if !got.OK || got.Command != "make-batch" {
		t.Fatalf("unexpected batch response: %+v", got)
	}
	if got.Summary.JobCount != 2 || got.Summary.SuccessCount != 2 || got.Summary.FailedCount != 0 {
		t.Fatalf("unexpected batch summary: %+v", got.Summary)
	}
	if len(got.Jobs) != 2 {
		t.Fatalf("expected two jobs, got %+v", got.Jobs)
	}
	if got.Jobs[0].ID != "job-a" || got.Jobs[1].ID != "job-b" {
		t.Fatalf("expected manifest-order jobs, got %+v", got.Jobs)
	}
	if got.Jobs[0].Outputs.PNG == nil || got.Jobs[1].Outputs.PNG == nil {
		t.Fatalf("expected png outputs for both jobs, got %+v", got.Jobs)
	}
}

func TestMakeBatchJSONManifestFailureIsExit4(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, manifest := writeBatchManifestIssueFixture(t, "{not-json}")

	code := run([]string{"make", "batch", manifest, "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 4 {
		t.Fatalf("expected exit code 4, got %d", code)
	}

	got := decodeMakeBatchResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "INVALID_BATCH_MANIFEST_JSON" {
		t.Fatalf("unexpected manifest error response: %+v", got)
	}
}

func TestMakeBatchJSONMissingManifestAndUnknownArgAreCLIUsage(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets := fixturePath("basic-assets")

	code := run([]string{"make", "batch", "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2 for missing manifest, got %d", code)
	}
	missingManifest := decodeMakeBatchResponse(t, stdout.Bytes())
	if missingManifest.OK || len(missingManifest.Errors) != 1 || missingManifest.Errors[0].Code != "MISSING_MANIFEST" {
		t.Fatalf("unexpected missing manifest response: %+v", missingManifest)
	}

	stdout.Reset()
	stderr.Reset()
	code = run([]string{"make", "batch", "manifest.json", "--assets", assets, "--out", "ignored.png", "--json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("expected exit code 2 for unknown argument, got %d", code)
	}
	unknownArg := decodeMakeBatchResponse(t, stdout.Bytes())
	if unknownArg.OK || len(unknownArg.Errors) != 1 || unknownArg.Errors[0].Code != "UNKNOWN_ARGUMENT" {
		t.Fatalf("unexpected unknown arg response: %+v", unknownArg)
	}
}

func TestMakeBatchJSONJobFailureMapsToExit3AndIncludesFailingContext(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, manifest := writeFailingBatchManifestFixture(t)

	code := run([]string{"make", "batch", manifest, "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("expected exit code 3, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}

	got := decodeMakeBatchResponse(t, stdout.Bytes())
	if got.OK || len(got.Jobs) != 1 {
		t.Fatalf("expected fail-fast with one job result, got %+v", got)
	}
	if got.Jobs[0].ID != "bad-first" || len(got.Jobs[0].Errors) != 1 {
		t.Fatalf("unexpected failing job payload: %+v", got.Jobs[0])
	}
	if got.Jobs[0].Errors[0].Code != "RECIPE_FILE_NOT_FOUND" {
		t.Fatalf("expected recipe failure code, got %+v", got.Jobs[0].Errors[0])
	}
	if len(got.Errors) != 1 {
		t.Fatalf("expected top-level failure, got %+v", got.Errors)
	}
	if got.Errors[0].Details["job_id"] != "bad-first" {
		t.Fatalf("expected failing job context in top-level error details, got %+v", got.Errors[0].Details)
	}
}

func TestMakeBatchJSONJobFailureMapsToExit6(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, manifest := writeRenderFailureBatchManifestFixture(t)

	code := run([]string{"make", "batch", manifest, "--assets", assets, "--json"}, &stdout, &stderr)
	if code != 6 {
		t.Fatalf("expected exit code 6, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}

	got := decodeMakeBatchResponse(t, stdout.Bytes())
	if got.OK || len(got.Errors) != 1 || got.Errors[0].Code != "RENDER_FAILED" {
		t.Fatalf("unexpected render failure response: %+v", got)
	}
}

func TestMakeBatchTextSuccessSummaryIsDeterministic(t *testing.T) {
	var stdout, stderr bytes.Buffer
	assets, manifest := writeSuccessfulBatchManifestFixture(t)

	code := run([]string{"make", "batch", manifest, "--assets", assets}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}

	want := "ok: make-batch\njobs: 2\nsuccess: 2\nfailed: 0\njob 1: job-a ok\njob 2: job-b ok\n"
	if stdout.String() != want {
		t.Fatalf("unexpected summary output\nwant:\n%s\ngot:\n%s", want, stdout.String())
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

func decodeAssetsValidationResponse(t *testing.T, data []byte) assetsValidationTestResponse {
	t.Helper()

	var got assetsValidationTestResponse
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

func writeMinimalAssetsPack(t *testing.T, assets string) {
	t.Helper()

	writeTestFile(t, filepath.Join(assets, "pack.json"), `{"schema_version":"1","id":"minimal-assets","name":"Minimal Assets"}`)
	writeTestFile(t, filepath.Join(assets, "sheet_definitions", "body", "body_human.json"), `{"name":"Human Body","type_name":"body","layer_1":{"zPos":10,"male":"body/human/male/"},"animations":["walk"]}`)
	writeTestFile(t, filepath.Join(assets, "spritesheets", ".gitkeep"), "")
	writeTestFile(t, filepath.Join(assets, "palette_definitions", ".gitkeep"), "")
}

func decodeMakeResponse(t *testing.T, data []byte) makeTestResponse {
	t.Helper()
	var got makeTestResponse
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to decode JSON %q: %v", string(data), err)
	}
	return got
}

func decodeMakeBatchResponse(t *testing.T, data []byte) makeBatchTestResponse {
	t.Helper()
	var got makeBatchTestResponse
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to decode JSON %q: %v", string(data), err)
	}
	return got
}

func decodeDownloadAssetsResponse(t *testing.T, data []byte) downloadAssetsTestResponse {
	t.Helper()
	var got downloadAssetsTestResponse
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("failed to decode JSON %q: %v", string(data), err)
	}
	return got
}

func assertHasKeys(t *testing.T, got map[string]interface{}, keys ...string) {
	t.Helper()
	for _, key := range keys {
		if _, ok := got[key]; !ok {
			t.Fatalf("missing key %q in %+v", key, got)
		}
	}
}

func writeDeterministicMakeFixture(t *testing.T) (string, string) {
	t.Helper()

	root := t.TempDir()
	assets := filepath.Join(root, "assets")
	writeTestFile(t, filepath.Join(assets, "pack.json"), `{"schema_version":"1","id":"make-assets","name":"Make Assets","defaults":{"body_type":"male","animations":["idle","walk"],"canvas_width":8}}`)
	writeTestFile(t, filepath.Join(assets, "sheet_definitions", "body", "body_human.json"), `{"name":"Human Body","type_name":"body","layer_1":{"zPos":10,"male":"body/human/male/"},"animations":["walk","idle"]}`)
	writeTestFile(t, filepath.Join(assets, "sheet_definitions", "weapon", "sword_training.json"), `{"name":"Training Sword","type_name":"weapon","layer_1":{"zPos":30,"male":"weapon/sword/male/"},"animations":["walk","idle"]}`)
	writeTestFile(t, filepath.Join(assets, "palette_definitions", ".gitkeep"), "")
	writeTestFile(t, filepath.Join(assets, "spritesheets", ".gitkeep"), "")

	makeDeterministicLayerPNG(t, filepath.Join(assets, "spritesheets", "body", "human", "male", "idle.png"), color.RGBA{R: 40, G: 100, B: 200, A: 255})
	makeDeterministicLayerPNG(t, filepath.Join(assets, "spritesheets", "body", "human", "male", "walk.png"), color.RGBA{R: 50, G: 110, B: 210, A: 255})
	makeDeterministicLayerPNG(t, filepath.Join(assets, "spritesheets", "weapon", "sword", "male", "idle.png"), color.RGBA{R: 200, G: 40, B: 80, A: 255})
	makeDeterministicLayerPNG(t, filepath.Join(assets, "spritesheets", "weapon", "sword", "male", "walk.png"), color.RGBA{R: 210, G: 50, B: 90, A: 255})

	recipe := filepath.Join(root, "recipe.json")
	writeTestFile(t, recipe, `{"body_type":"male","selections":{"body":{"id":"body_human"},"weapon":{"id":"sword_training"}}}`)
	return assets, recipe
}

func makeDeterministicLayerPNG(t *testing.T, path string, fill color.RGBA) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			c := fill
			c.R = uint8((int(c.R) + x + y) % 255)
			img.SetRGBA(x, y, c)
		}
	}
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		t.Fatal(err)
	}
}

type readinessCLIFixtureOptions struct {
	recipeBodyType        string
	missingFallback       string
	requiredAnimations    []string
	includeBodyIdle       bool
	includeBodyFemaleIdle bool
	includeBodyWalk       bool
	includeWeaponIdle     bool
	includeWeaponWalk     bool
	weaponHasFemalePath   bool
}

func writeReadinessCLIFixture(t *testing.T, opts readinessCLIFixtureOptions) (string, string) {
	t.Helper()
	root := t.TempDir()
	assets := filepath.Join(root, "assets")

	recipeBodyType := opts.recipeBodyType
	if recipeBodyType == "" {
		recipeBodyType = "male"
	}
	missingFallback := opts.missingFallback
	if missingFallback == "" {
		missingFallback = "male"
	}
	requiredAnimations := opts.requiredAnimations
	if len(requiredAnimations) == 0 {
		requiredAnimations = []string{"idle"}
	}
	animationsJSON := `["` + requiredAnimations[0] + `"`
	for i := 1; i < len(requiredAnimations); i++ {
		animationsJSON += `,"` + requiredAnimations[i] + `"`
	}
	animationsJSON += `]`

	weaponFemale := `"female":"weapon/sword/female/"`
	if !opts.weaponHasFemalePath {
		weaponFemale = ""
	}

	writeTestFile(t, filepath.Join(assets, "pack.json"), `{"schema_version":"1","id":"cli-readiness","name":"CLI Readiness","defaults":{"body_type":"male","animations":`+animationsJSON+`,"canvas_width":8,"missing_body_type_fallback":"`+missingFallback+`"}}`)
	writeTestFile(t, filepath.Join(assets, "sheet_definitions", "body", "body_human.json"), `{"name":"Human Body","type_name":"body","layer_1":{"zPos":10,"male":"body/human/male/","female":"body/human/female/"},"animations":["idle","walk"]}`)
	writeTestFile(t, filepath.Join(assets, "sheet_definitions", "weapon", "sword_training.json"), `{"name":"Training Sword","type_name":"weapon","layer_1":{"zPos":30,"male":"weapon/sword/male/"`+leadingCommaIfNotEmpty(weaponFemale)+`},"animations":["idle","walk"]}`)
	writeTestFile(t, filepath.Join(assets, "palette_definitions", ".gitkeep"), "")
	writeTestFile(t, filepath.Join(assets, "spritesheets", ".gitkeep"), "")

	if opts.includeBodyIdle {
		makeDeterministicLayerPNG(t, filepath.Join(assets, "spritesheets", "body", "human", "male", "idle.png"), color.RGBA{R: 40, G: 100, B: 200, A: 255})
	}
	if opts.includeBodyWalk {
		makeDeterministicLayerPNG(t, filepath.Join(assets, "spritesheets", "body", "human", "male", "walk.png"), color.RGBA{R: 50, G: 110, B: 210, A: 255})
	}
	if opts.includeBodyFemaleIdle {
		makeDeterministicLayerPNG(t, filepath.Join(assets, "spritesheets", "body", "human", "female", "idle.png"), color.RGBA{R: 44, G: 104, B: 206, A: 255})
	}
	if opts.includeWeaponIdle {
		makeDeterministicLayerPNG(t, filepath.Join(assets, "spritesheets", "weapon", "sword", "male", "idle.png"), color.RGBA{R: 200, G: 40, B: 80, A: 255})
	}
	if opts.includeWeaponWalk {
		makeDeterministicLayerPNG(t, filepath.Join(assets, "spritesheets", "weapon", "sword", "male", "walk.png"), color.RGBA{R: 210, G: 50, B: 90, A: 255})
	}

	recipe := filepath.Join(root, "recipe.json")
	writeTestFile(t, recipe, `{"body_type":"`+recipeBodyType+`","selections":{"body":{"id":"body_human"},"weapon":{"id":"sword_training"}}}`)
	return assets, recipe
}

func leadingCommaIfNotEmpty(value string) string {
	if value == "" {
		return ""
	}
	return "," + value
}

func pngSHA256(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	return fmt.Sprintf("%x", sum[:])
}

func writeSuccessfulBatchManifestFixture(t *testing.T) (string, string) {
	t.Helper()

	assets, recipe := writeDeterministicMakeFixture(t)
	root := t.TempDir()
	manifest := filepath.Join(root, "batch.json")
	manifestText := fmt.Sprintf(`{
  "schema_version": "1",
  "jobs": [
    {
      "id": "job-a",
      "recipe": %q,
      "out": %q,
      "report": %q
    },
    {
      "id": "job-b",
      "recipe": %q,
      "out": %q
    }
  ]
}`, recipe, filepath.Join(root, "out", "job-a.png"), filepath.Join(root, "reports", "job-a.report.json"), recipe, filepath.Join(root, "out", "job-b.png"))
	writeTestFile(t, manifest, manifestText)
	return assets, manifest
}

func writeBatchManifestIssueFixture(t *testing.T, payload string) (string, string) {
	t.Helper()
	assets := fixturePath("basic-assets")
	manifest := filepath.Join(t.TempDir(), "batch.json")
	writeTestFile(t, manifest, payload)
	return assets, manifest
}

func writeFailingBatchManifestFixture(t *testing.T) (string, string) {
	t.Helper()

	assets, recipe := writeDeterministicMakeFixture(t)
	root := t.TempDir()
	manifest := filepath.Join(root, "batch.json")
	manifestText := fmt.Sprintf(`{
  "schema_version": "1",
  "jobs": [
    {
      "id": "bad-first",
      "recipe": %q,
      "out": %q
    },
    {
      "id": "would-pass",
      "recipe": %q,
      "out": %q
    }
  ]
}`, filepath.Join(root, "missing.json"), filepath.Join(root, "out", "bad.png"), recipe, filepath.Join(root, "out", "good.png"))
	writeTestFile(t, manifest, manifestText)
	return assets, manifest
}

func writeRenderFailureBatchManifestFixture(t *testing.T) (string, string) {
	t.Helper()

	assets, recipe := writeDeterministicMakeFixture(t)
	root := t.TempDir()
	outDir := filepath.Join(root, "out-dir")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := filepath.Join(root, "batch.json")
	manifestText := fmt.Sprintf(`{
  "schema_version": "1",
  "jobs": [
    {
      "id": "render-fail",
      "recipe": %q,
      "out": %q
    }
  ]
}`, recipe, outDir)
	writeTestFile(t, manifest, manifestText)
	return assets, manifest
}

func writeSlashParityCLIFixture(t *testing.T, includeMappedSlash bool) (string, string) {
	t.Helper()
	root := t.TempDir()
	assets := filepath.Join(root, "assets")

	writeTestFile(t, filepath.Join(assets, "pack.json"), `{"schema_version":"1","id":"cli-slash-parity","name":"CLI Slash Parity","defaults":{"body_type":"male","animations":["slash"],"canvas_width":8}}`)
	writeTestFile(t, filepath.Join(assets, "sheet_definitions", "body", "body_human.json"), `{"name":"Human Body","type_name":"body","layer_1":{"zPos":10,"male":"body/human/male/"},"animations":["slash"]}`)
	writeTestFile(t, filepath.Join(assets, "sheet_definitions", "weapon", "sword_training.json"), `{"name":"Training Sword","type_name":"weapon","layer_1":{"zPos":30,"male":"weapon/sword/male/"},"animations":["slash"]}`)
	writeTestFile(t, filepath.Join(assets, "palette_definitions", ".gitkeep"), "")
	writeTestFile(t, filepath.Join(assets, "spritesheets", ".gitkeep"), "")

	makeDeterministicLayerPNG(t, filepath.Join(assets, "spritesheets", "body", "human", "male", "slash.png"), color.RGBA{R: 40, G: 100, B: 200, A: 255})
	if includeMappedSlash {
		makeDeterministicLayerPNG(t, filepath.Join(assets, "spritesheets", "weapon", "sword", "male", "attack_slash", "front.png"), color.RGBA{R: 200, G: 40, B: 80, A: 255})
		makeDeterministicLayerPNG(t, filepath.Join(assets, "spritesheets", "weapon", "sword", "male", "attack_slash", "pose_behind.png"), color.RGBA{R: 180, G: 20, B: 60, A: 255})
	}

	recipe := filepath.Join(root, "recipe.json")
	writeTestFile(t, recipe, `{"body_type":"male","selections":{"body":{"id":"body_human"},"weapon":{"id":"sword_training"}}}`)
	return assets, recipe
}

func createLPCArchiveZipPayload(t *testing.T, includeInvalidPack bool) []byte {
	t.Helper()
	buffer := &bytes.Buffer{}
	archive := zip.NewWriter(buffer)

	base := "Universal-LPC-Spritesheet-Character-Generator-master/"
	addZipFile := func(path string, contents string) {
		writer, err := archive.Create(base + strings.TrimPrefix(filepath.ToSlash(path), "/"))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.WriteString(writer, contents); err != nil {
			t.Fatal(err)
		}
	}

	if includeInvalidPack {
		addZipFile("pack.json", `{not-json}`)
	}
	addZipFile("sheet_definitions/body/body_human.json", `{"name":"Human Body","type_name":"body","layer_1":{"zPos":10,"male":"body/human/male/"},"animations":["walk"]}`)
	addZipFile("spritesheets/.gitkeep", "")
	addZipFile("palette_definitions/.gitkeep", "")

	if err := archive.Close(); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}
