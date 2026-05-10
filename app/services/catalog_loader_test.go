package services

import (
	"path/filepath"
	"testing"
)

func TestCatalogLoaderLoadsPackAndGroupsLayers(t *testing.T) {
	assets := filepath.Join("..", "..", "testdata", "fixtures", "basic-assets")

	catalog, loadErr := NewCatalogLoader().Load(assets)
	if loadErr != nil {
		t.Fatalf("expected catalog to load, got %v", loadErr)
	}

	if catalog.Pack.ID != "basic-test-assets" {
		t.Fatalf("unexpected pack id: %q", catalog.Pack.ID)
	}
	if len(catalog.Categories) != 2 {
		t.Fatalf("expected 2 categories, got %+v", catalog.Categories)
	}
	if catalog.Categories[0].ID != "body" || catalog.Categories[1].ID != "weapon" {
		t.Fatalf("categories should be sorted by id, got %+v", catalog.Categories)
	}

	bodyLayer := catalog.Categories[0].Layers[0]
	if bodyLayer.ID != "body_human" {
		t.Fatalf("unexpected layer id: %q", bodyLayer.ID)
	}
	if bodyLayer.PathPrefix != "body/human/female/" {
		t.Fatalf("expected deterministic first body-type path, got %q", bodyLayer.PathPrefix)
	}
}

func TestCatalogLoaderMissingPackJSON(t *testing.T) {
	_, loadErr := NewCatalogLoader().Load(t.TempDir())
	if loadErr == nil {
		t.Fatal("expected missing pack error")
	}
	if loadErr.Code != "MISSING_PACK_JSON" {
		t.Fatalf("unexpected error: %+v", loadErr)
	}
}

func TestCatalogLoaderAssetsDirectoryNotFound(t *testing.T) {
	_, loadErr := NewCatalogLoader().Load(filepath.Join(t.TempDir(), "missing-assets"))
	if loadErr == nil {
		t.Fatal("expected missing assets directory error")
	}
	if loadErr.Code != "ASSETS_DIRECTORY_NOT_FOUND" {
		t.Fatalf("unexpected error: %+v", loadErr)
	}
}
