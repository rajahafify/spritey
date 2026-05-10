package services

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/rajahafify/spritey/app/models"
)

type RecipeValidator struct {
	loader CatalogLoader
}

func NewRecipeValidator() RecipeValidator {
	return RecipeValidator{loader: NewCatalogLoader()}
}

func (validator RecipeValidator) Validate(recipePath string, assetsPath string) (models.RecipeValidationResult, *models.Problem) {
	catalog, problem := validator.loader.Load(assetsPath)
	if problem != nil {
		return models.RecipeValidationResult{}, problem
	}

	recipe, problem := loadRecipe(recipePath)
	if problem != nil {
		return models.RecipeValidationResult{}, problem
	}

	bodyType := recipe.BodyType
	if bodyType == "" {
		bodyType = catalog.Pack.Defaults.BodyType
	}

	if len(recipe.Selections) == 0 {
		return models.RecipeValidationResult{}, &models.Problem{
			Code:    "MISSING_SELECTIONS",
			Message: "recipe must contain at least one selection",
			Field:   "selections",
		}
	}

	layerIndex := indexLayers(catalog)
	selectionKeys := sortedSelectionKeys(recipe.Selections)
	validated := make([]models.RecipeValidationSelection, 0, len(selectionKeys))
	for _, category := range selectionKeys {
		selection := recipe.Selections[category]
		if selection.ID == "" {
			return models.RecipeValidationResult{}, &models.Problem{
				Code:    "MISSING_SELECTION_ID",
				Message: fmt.Sprintf("selection %s must contain an id", category),
				Field:   fmt.Sprintf("selections.%s.id", category),
			}
		}

		layer, ok := layerIndex[selection.ID]
		if !ok {
			return models.RecipeValidationResult{}, &models.Problem{
				Code:    "UNKNOWN_LAYER_ID",
				Message: fmt.Sprintf("layer id not found: %s", selection.ID),
				Field:   fmt.Sprintf("selections.%s.id", category),
			}
		}
		if !supportsBodyType(layer.Layer.BodyTypes, bodyType) {
			return models.RecipeValidationResult{}, &models.Problem{
				Code:    "UNSUPPORTED_BODY_TYPE",
				Message: fmt.Sprintf("layer %s does not support body type %s", selection.ID, bodyType),
				Field:   "body_type",
			}
		}

		validated = append(validated, models.RecipeValidationSelection{
			Category:       category,
			ID:             selection.ID,
			PaletteVariant: selection.PaletteVariant,
		})
	}

	return models.RecipeValidationResult{
		Path:       recipePath,
		BodyType:   bodyType,
		Selections: validated,
	}, nil
}

func loadRecipe(path string) (models.Recipe, *models.Problem) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return models.Recipe{}, &models.Problem{
				Code:    "RECIPE_FILE_NOT_FOUND",
				Message: "recipe file does not exist",
				Field:   "recipe",
			}
		}
		return models.Recipe{}, &models.Problem{
			Code:    "READ_RECIPE_FAILED",
			Message: err.Error(),
			Field:   "recipe",
		}
	}

	var recipe models.Recipe
	if err := json.Unmarshal(data, &recipe); err != nil {
		return models.Recipe{}, &models.Problem{
			Code:    "INVALID_RECIPE_JSON",
			Message: fmt.Sprintf("recipe is not valid JSON: %v", err),
			Field:   "recipe",
		}
	}
	return recipe, nil
}

type indexedLayer struct {
	Category string
	Layer    models.Layer
}

func indexLayers(catalog models.Catalog) map[string]indexedLayer {
	index := map[string]indexedLayer{}
	for _, category := range catalog.Categories {
		for _, layer := range category.Layers {
			if _, exists := index[layer.ID]; exists {
				continue
			}
			index[layer.ID] = indexedLayer{Category: category.ID, Layer: layer}
		}
	}
	return index
}

func sortedSelectionKeys(selections map[string]models.RecipeSelection) []string {
	keys := make([]string, 0, len(selections))
	for key := range selections {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func supportsBodyType(bodyTypes []string, bodyType string) bool {
	for _, candidate := range bodyTypes {
		if candidate == bodyType {
			return true
		}
	}
	return false
}
