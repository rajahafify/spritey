package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/rajahafify/spritey/app/models"
)

type RecipeValidator struct {
	loader        CatalogLoader
	frameResolver SpriteFrameResolver
}

func NewRecipeValidator() RecipeValidator {
	return RecipeValidator{
		loader:        NewCatalogLoader(),
		frameResolver: NewSpriteFrameResolver(),
	}
}

func (validator RecipeValidator) Validate(recipePath string, assetsPath string) (models.RecipeValidationResult, *models.Problem) {
	return validator.validate(recipePath, assetsPath, true)
}

func (validator RecipeValidator) ValidateForMake(recipePath string, assetsPath string) (models.RecipeValidationResult, *models.Problem) {
	return validator.validate(recipePath, assetsPath, false)
}

func (validator RecipeValidator) validate(recipePath string, assetsPath string, requireFrames bool) (models.RecipeValidationResult, *models.Problem) {
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
	requiredAnimations := requiredAnimationIDs(catalog.Pack.Defaults.Animations)

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
	renderInputs := make([]models.RecipeRenderInput, 0, len(selectionKeys))
	warnings := []string{}
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

		resolvedPath, resolvedBodyType, usedFallback, ok := resolveLayerRenderPath(layer.Layer, bodyType, catalog.Pack.Defaults.MissingBodyTypeFallback)
		if !ok {
			return models.RecipeValidationResult{}, &models.Problem{
				Code:    "UNSUPPORTED_BODY_TYPE",
				Message: fmt.Sprintf("layer %s does not support body type %s", selection.ID, bodyType),
				Field:   "body_type",
			}
		}
		if usedFallback {
			warnings = append(warnings, fallbackBodyTypeWarning(selection.ID, bodyType, resolvedBodyType, resolvedPath))
		}
		if requireFrames {
			for _, animationID := range requiredAnimations {
				_, found, err := validator.frameResolver.ResolveFrame(assetsPath, resolvedPath, animationID)
				if err != nil {
					return models.RecipeValidationResult{}, &models.Problem{
						Code:    "READ_SPRITE_FRAME_FAILED",
						Message: err.Error(),
						Field:   path.Join("spritesheets", path.Clean(path.Join(resolvedPath, animationID+".png"))),
					}
				}
				if !found {
					relativeFrame := path.Join("spritesheets", path.Clean(path.Join(resolvedPath, animationID+".png")))
					return models.RecipeValidationResult{}, &models.Problem{
						Code:    "MISSING_SPRITE_FRAME",
						Message: fmt.Sprintf("missing required sprite frame: %s", relativeFrame),
						Field:   relativeFrame,
					}
				}
			}
		}

		validated = append(validated, models.RecipeValidationSelection{
			Category:       category,
			ID:             selection.ID,
			PaletteVariant: selection.PaletteVariant,
		})
		renderInputs = append(renderInputs, models.RecipeRenderInput{
			Category:         category,
			LayerID:          selection.ID,
			ResolvedPath:     resolvedPath,
			ResolvedBodyType: resolvedBodyType,
		})
	}

	return models.RecipeValidationResult{
		Path:                 recipePath,
		BodyTypeRequested:    recipe.BodyType,
		BodyType:             bodyType,
		Selections:           validated,
		Warnings:             warnings,
		RequiredAnimationIDs: requiredAnimations,
		RenderInputs:         renderInputs,
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

func requiredAnimationIDs(animations []string) []string {
	if len(animations) == 0 {
		return []string{"idle"}
	}

	cleaned := make([]string, 0, len(animations))
	for _, animation := range animations {
		trimmed := strings.TrimSpace(animation)
		if trimmed == "" {
			continue
		}
		cleaned = append(cleaned, trimmed)
	}
	if len(cleaned) == 0 {
		return []string{"idle"}
	}
	sort.Strings(cleaned)
	unique := cleaned[:1]
	for i := 1; i < len(cleaned); i++ {
		if cleaned[i] != cleaned[i-1] {
			unique = append(unique, cleaned[i])
		}
	}
	return unique
}

func resolveLayerRenderPath(layer models.Layer, requestedBodyType string, fallbackBodyType string) (string, string, bool, bool) {
	if layer.BodyTypePaths != nil {
		if pathPrefix := strings.TrimSpace(layer.BodyTypePaths[requestedBodyType]); pathPrefix != "" {
			return pathPrefix, requestedBodyType, false, true
		}
	}

	effectiveFallback := strings.TrimSpace(fallbackBodyType)
	if effectiveFallback != "" && effectiveFallback != requestedBodyType && layer.BodyTypePaths != nil {
		if pathPrefix := strings.TrimSpace(layer.BodyTypePaths[effectiveFallback]); pathPrefix != "" {
			return pathPrefix, effectiveFallback, true, true
		}
	}

	return "", "", false, false
}

func fallbackBodyTypeWarning(layerID string, requestedBodyType string, fallbackBodyType string, resolvedPath string) string {
	return fmt.Sprintf(
		"layer %s missing body type %s; using fallback %s path %s",
		layerID,
		requestedBodyType,
		fallbackBodyType,
		resolvedPath,
	)
}
