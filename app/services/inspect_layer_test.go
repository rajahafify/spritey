package services

import (
	"path/filepath"
	"testing"
)

func TestInspectLayerFindsExactLayerAcrossCategories(t *testing.T) {
	assets := filepath.Join("..", "..", "testdata", "fixtures", "basic-assets")

	result, problem := NewInspectLayerService().Find(assets, "sword_training")
	if problem != nil {
		t.Fatalf("expected layer to load, got %v", problem)
	}

	if result.Category != "weapon" {
		t.Fatalf("expected weapon category, got %q", result.Category)
	}
	if result.ID != "sword_training" {
		t.Fatalf("unexpected layer: %+v", result)
	}
}

func TestInspectLayerUnknownLayerID(t *testing.T) {
	assets := filepath.Join("..", "..", "testdata", "fixtures", "basic-assets")

	_, problem := NewInspectLayerService().Find(assets, "missing_layer")
	if problem == nil {
		t.Fatal("expected unknown layer problem")
	}
	if problem.Code != "UNKNOWN_LAYER_ID" {
		t.Fatalf("unexpected problem: %+v", problem)
	}
}
