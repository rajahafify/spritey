package services

import (
	"path/filepath"
	"testing"
)

func TestAssetsValidatorBuildsSummaryCounts(t *testing.T) {
	assets := filepath.Join("..", "..", "testdata", "fixtures", "basic-assets")

	result, problem := NewAssetsValidator().Validate(assets)
	if problem != nil {
		t.Fatalf("expected assets validation to pass, got %v", problem)
	}

	if result.Assets.Path != assets {
		t.Fatalf("unexpected assets path: %q", result.Assets.Path)
	}
	if result.Pack.ID != "basic-test-assets" {
		t.Fatalf("unexpected pack id: %q", result.Pack.ID)
	}
	if result.Summary.CategoryCount != 2 {
		t.Fatalf("expected 2 categories, got %d", result.Summary.CategoryCount)
	}
	if result.Summary.LayerCount != 2 {
		t.Fatalf("expected 2 layers, got %d", result.Summary.LayerCount)
	}
}
