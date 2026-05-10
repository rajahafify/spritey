package services

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rajahafify/spritey/app/models"
)

type CatalogLoader struct{}

func NewCatalogLoader() CatalogLoader {
	return CatalogLoader{}
}

func (CatalogLoader) Load(assetsPath string) (models.Catalog, *models.Problem) {
	if info, err := os.Stat(assetsPath); err != nil || !info.IsDir() {
		if os.IsNotExist(err) {
			return models.Catalog{}, &models.Problem{
				Code:    "ASSETS_DIRECTORY_NOT_FOUND",
				Message: "assets directory does not exist",
				Field:   "assets",
			}
		}
		if err != nil {
			return models.Catalog{}, &models.Problem{
				Code:    "READ_ASSETS_DIRECTORY_FAILED",
				Message: err.Error(),
				Field:   "assets",
			}
		}
		return models.Catalog{}, &models.Problem{
			Code:    "ASSETS_PATH_NOT_DIRECTORY",
			Message: "assets path must be a directory",
			Field:   "assets",
		}
	}

	pack, problem := loadPack(filepath.Join(assetsPath, "pack.json"))
	if problem != nil {
		return models.Catalog{}, problem
	}

	sheetRoot := filepath.Join(assetsPath, "sheet_definitions")
	if info, err := os.Stat(sheetRoot); err != nil || !info.IsDir() {
		return models.Catalog{}, &models.Problem{
			Code:    "MISSING_SHEET_DEFINITIONS",
			Message: "sheet_definitions directory is required",
			Field:   "sheet_definitions",
		}
	}

	categoryMap := map[string][]models.Layer{}
	walkErr := filepath.WalkDir(sheetRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if filepath.Ext(entry.Name()) != ".json" || strings.HasPrefix(entry.Name(), "meta_") {
			return nil
		}

		layer, categoryID, problem := loadLayerDefinition(sheetRoot, path)
		if problem != nil {
			return problem
		}
		categoryMap[categoryID] = append(categoryMap[categoryID], layer)
		return nil
	})
	if walkErr != nil {
		if problem, ok := walkErr.(*models.Problem); ok {
			return models.Catalog{}, problem
		}
		return models.Catalog{}, &models.Problem{
			Code:    "READ_SHEET_DEFINITIONS_FAILED",
			Message: walkErr.Error(),
			Field:   "sheet_definitions",
		}
	}

	return models.Catalog{
		Pack:       pack,
		Categories: sortedCategories(categoryMap),
		Warnings:   []string{},
	}, nil
}

func loadPack(path string) (models.Pack, *models.Problem) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return models.Pack{}, &models.Problem{
				Code:    "MISSING_PACK_JSON",
				Message: "pack.json is required in the assets directory",
				Field:   "pack.json",
			}
		}
		return models.Pack{}, &models.Problem{
			Code:    "READ_PACK_JSON_FAILED",
			Message: err.Error(),
			Field:   "pack.json",
		}
	}

	var pack models.Pack
	if err := json.Unmarshal(data, &pack); err != nil {
		return models.Pack{}, &models.Problem{
			Code:    "INVALID_PACK_JSON",
			Message: fmt.Sprintf("pack.json is not valid JSON: %v", err),
			Field:   "pack.json",
		}
	}
	return pack, nil
}

type rawSheetDefinition struct {
	Name       string                 `json:"name"`
	TypeName   string                 `json:"type_name"`
	Layer1     map[string]interface{} `json:"layer_1"`
	Animations []string               `json:"animations"`
	Recolors   struct {
		Material string `json:"material"`
	} `json:"recolors"`
	Credits []models.Credit `json:"credits"`
}

func loadLayerDefinition(sheetRoot string, path string) (models.Layer, string, *models.Problem) {
	data, err := os.ReadFile(path)
	field := slashRel(sheetRoot, path)
	if err != nil {
		return models.Layer{}, "", &models.Problem{
			Code:    "READ_SHEET_DEFINITION_FAILED",
			Message: err.Error(),
			Field:   field,
		}
	}

	var definition rawSheetDefinition
	if err := json.Unmarshal(data, &definition); err != nil {
		return models.Layer{}, "", &models.Problem{
			Code:    "INVALID_SHEET_DEFINITION_JSON",
			Message: fmt.Sprintf("sheet definition is not valid JSON: %v", err),
			Field:   field,
		}
	}

	id := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	name := definition.Name
	if name == "" {
		name = id
	}

	bodyTypes, bodyTypePaths, pathPrefix := layerBodyTypes(definition.Layer1)
	credits := definition.Credits
	if credits == nil {
		credits = []models.Credit{}
	}
	animations := definition.Animations
	if animations == nil {
		animations = []string{}
	}

	layer := models.Layer{
		ID:              id,
		Name:            name,
		ZPos:            layerZPos(definition.Layer1),
		BodyTypes:       bodyTypes,
		BodyTypePaths:   bodyTypePaths,
		Animations:      animations,
		RecolorMaterial: definition.Recolors.Material,
		PathPrefix:      pathPrefix,
		Credits:         credits,
	}

	return layer, categoryID(sheetRoot, path, definition.TypeName), nil
}

func layerZPos(layer map[string]interface{}) int {
	if layer == nil {
		return 0
	}
	switch value := layer["zPos"].(type) {
	case float64:
		return int(value)
	case int:
		return value
	default:
		return 0
	}
}

func layerBodyTypes(layer map[string]interface{}) ([]string, map[string]string, string) {
	if layer == nil {
		return []string{}, map[string]string{}, ""
	}

	bodyTypes := make([]string, 0, len(layer))
	pathsByBodyType := map[string]string{}
	for key, value := range layer {
		if isLayerMetadataKey(key) {
			continue
		}
		pathPrefix, ok := value.(string)
		if !ok || pathPrefix == "" {
			continue
		}
		bodyTypes = append(bodyTypes, key)
		pathsByBodyType[key] = pathPrefix
	}
	sort.Strings(bodyTypes)
	if len(bodyTypes) == 0 {
		return bodyTypes, pathsByBodyType, ""
	}
	return bodyTypes, pathsByBodyType, pathsByBodyType[bodyTypes[0]]
}

func isLayerMetadataKey(key string) bool {
	switch key {
	case "zPos", "custom_animation":
		return true
	default:
		return false
	}
}

func categoryID(sheetRoot string, path string, explicitTypeName string) string {
	if explicitTypeName != "" {
		return explicitTypeName
	}

	rel := slashRel(sheetRoot, filepath.Dir(path))
	if rel == "." || rel == "" {
		return "uncategorized"
	}
	return strings.Split(rel, "/")[0]
}

func sortedCategories(categoryMap map[string][]models.Layer) []models.Category {
	ids := make([]string, 0, len(categoryMap))
	for id := range categoryMap {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	categories := make([]models.Category, 0, len(ids))
	for _, id := range ids {
		layers := categoryMap[id]
		sort.Slice(layers, func(i int, j int) bool {
			return layers[i].ID < layers[j].ID
		})
		categories = append(categories, models.Category{
			ID:     id,
			Layers: layers,
		})
	}
	return categories
}

func slashRel(base string, path string) string {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(rel)
}
