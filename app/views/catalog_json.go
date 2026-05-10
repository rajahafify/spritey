package views

import (
	"encoding/json"
	"io"

	"github.com/rajahafify/spritey/app/models"
)

type CatalogResponse struct {
	OK         bool              `json:"ok"`
	Pack       *models.Pack      `json:"pack"`
	Categories []models.Category `json:"categories"`
	Warnings   []string          `json:"warnings"`
	Errors     []models.Problem  `json:"errors"`
}

func WriteCatalogJSON(writer io.Writer, catalog models.Catalog) error {
	return writeJSON(writer, CatalogResponse{
		OK:         true,
		Pack:       &catalog.Pack,
		Categories: catalog.Categories,
		Warnings:   catalog.Warnings,
		Errors:     []models.Problem{},
	})
}

func WriteCatalogErrorJSON(writer io.Writer, problem models.Problem) error {
	return writeJSON(writer, CatalogResponse{
		OK:         false,
		Pack:       nil,
		Categories: []models.Category{},
		Warnings:   []string{},
		Errors:     []models.Problem{problem},
	})
}

func writeJSON(writer io.Writer, value interface{}) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
