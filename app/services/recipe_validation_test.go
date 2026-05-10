package services

import (
	"path/filepath"
	"testing"
)

func TestRecipeValidatorValidatesRecipe(t *testing.T) {
	assets := filepath.Join("..", "..", "testdata", "fixtures", "basic-assets")
	recipe := filepath.Join("..", "..", "testdata", "fixtures", "recipes", "valid-basic.json")

	result, problem := NewRecipeValidator().Validate(recipe, assets)
	if problem != nil {
		t.Fatalf("expected recipe to validate, got %v", problem)
	}

	if result.BodyType != "male" {
		t.Fatalf("unexpected body type: %q", result.BodyType)
	}
	if len(result.Selections) != 2 {
		t.Fatalf("expected 2 selections, got %+v", result.Selections)
	}
}

func TestRecipeValidatorUsesDefaultBodyType(t *testing.T) {
	assets := filepath.Join("..", "..", "testdata", "fixtures", "basic-assets")
	recipe := filepath.Join("..", "..", "testdata", "fixtures", "recipes", "default-body-type.json")

	result, problem := NewRecipeValidator().Validate(recipe, assets)
	if problem != nil {
		t.Fatalf("expected recipe to validate, got %v", problem)
	}
	if result.BodyType != "male" {
		t.Fatalf("expected default body type male, got %q", result.BodyType)
	}
}

func TestRecipeValidatorUnknownLayer(t *testing.T) {
	assets := filepath.Join("..", "..", "testdata", "fixtures", "basic-assets")
	recipe := filepath.Join("..", "..", "testdata", "fixtures", "recipes", "unknown-layer.json")

	_, problem := NewRecipeValidator().Validate(recipe, assets)
	if problem == nil {
		t.Fatal("expected validation problem")
	}
	if problem.Code != "UNKNOWN_LAYER_ID" {
		t.Fatalf("unexpected problem: %+v", problem)
	}
}
