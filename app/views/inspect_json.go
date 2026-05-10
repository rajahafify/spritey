package views

import (
	"io"

	"github.com/rajahafify/spritey/app/models"
)

type InspectLayerResponse struct {
	OK       bool                       `json:"ok"`
	Layer    *models.InspectLayerResult `json:"layer"`
	Warnings []string                   `json:"warnings"`
	Errors   []models.Problem           `json:"errors"`
}

func WriteInspectLayerJSON(writer io.Writer, layer models.InspectLayerResult) error {
	return writeJSON(writer, InspectLayerResponse{
		OK:       true,
		Layer:    &layer,
		Warnings: []string{},
		Errors:   []models.Problem{},
	})
}

func WriteInspectLayerErrorJSON(writer io.Writer, problem models.Problem) error {
	return writeJSON(writer, InspectLayerResponse{
		OK:       false,
		Layer:    nil,
		Warnings: []string{},
		Errors:   []models.Problem{problem},
	})
}
