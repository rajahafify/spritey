package services

import (
	"fmt"

	"github.com/rajahafify/spritey/app/models"
)

type InspectLayerService struct {
	loader CatalogLoader
}

func NewInspectLayerService() InspectLayerService {
	return InspectLayerService{loader: NewCatalogLoader()}
}

func (service InspectLayerService) Find(assetsPath string, layerID string) (models.InspectLayerResult, *models.Problem) {
	catalog, problem := service.loader.Load(assetsPath)
	if problem != nil {
		return models.InspectLayerResult{}, problem
	}

	for _, category := range catalog.Categories {
		for _, layer := range category.Layers {
			if layer.ID == layerID {
				return models.InspectLayerResult{
					Category:        category.ID,
					ID:              layer.ID,
					Name:            layer.Name,
					ZPos:            layer.ZPos,
					BodyTypes:       layer.BodyTypes,
					Animations:      layer.Animations,
					RecolorMaterial: layer.RecolorMaterial,
					PathPrefix:      layer.PathPrefix,
					Credits:         layer.Credits,
				}, nil
			}
		}
	}

	return models.InspectLayerResult{}, &models.Problem{
		Code:    "UNKNOWN_LAYER_ID",
		Message: fmt.Sprintf("layer id not found: %s", layerID),
		Field:   "layer_id",
	}
}
