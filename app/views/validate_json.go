package views

import (
	"io"

	"github.com/rajahafify/spritey/app/models"
)

type ValidateRecipeResponse struct {
	OK       bool                           `json:"ok"`
	Recipe   *models.RecipeValidationResult `json:"recipe"`
	Warnings []string                       `json:"warnings"`
	Errors   []models.Problem               `json:"errors"`
}

func WriteValidateRecipeJSON(writer io.Writer, result models.RecipeValidationResult) error {
	return writeJSON(writer, ValidateRecipeResponse{
		OK:       true,
		Recipe:   &result,
		Warnings: []string{},
		Errors:   []models.Problem{},
	})
}

func WriteValidateRecipeErrorJSON(writer io.Writer, problem models.Problem) error {
	return writeJSON(writer, ValidateRecipeResponse{
		OK:       false,
		Recipe:   nil,
		Warnings: []string{},
		Errors:   []models.Problem{problem},
	})
}
