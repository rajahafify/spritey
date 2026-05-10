package services

import (
	"os"
	"path/filepath"

	"github.com/rajahafify/spritey/app/models"
)

type AssetsValidator struct {
	loader CatalogLoader
}

func NewAssetsValidator() AssetsValidator {
	return AssetsValidator{loader: NewCatalogLoader()}
}

func (validator AssetsValidator) Validate(assetsPath string) (models.AssetsValidationResult, *models.Problem) {
	catalog, problem := validator.loader.Load(assetsPath)
	if problem != nil {
		return models.AssetsValidationResult{}, problem
	}

	if problem := requireAssetsSubdirectory(assetsPath, "spritesheets", "MISSING_SPRITESHEETS"); problem != nil {
		return models.AssetsValidationResult{}, problem
	}
	if problem := requireAssetsSubdirectory(assetsPath, "palette_definitions", "MISSING_PALETTE_DEFINITIONS"); problem != nil {
		return models.AssetsValidationResult{}, problem
	}

	return models.AssetsValidationResult{
		Assets: models.AssetsValidationTarget{
			Path: assetsPath,
		},
		Pack: catalog.Pack,
		Summary: models.AssetsValidationSummary{
			CategoryCount: len(catalog.Categories),
			LayerCount:    layerCount(catalog.Categories),
		},
		Warnings: catalog.Warnings,
	}, nil
}

func requireAssetsSubdirectory(assetsPath string, name string, code string) *models.Problem {
	info, err := os.Stat(filepath.Join(assetsPath, name))
	if err != nil || !info.IsDir() {
		return &models.Problem{
			Code:    code,
			Message: name + " directory is required",
			Field:   name,
		}
	}
	return nil
}

func layerCount(categories []models.Category) int {
	count := 0
	for _, category := range categories {
		count += len(category.Layers)
	}
	return count
}
