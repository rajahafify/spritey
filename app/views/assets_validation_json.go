package views

import (
	"io"

	"github.com/rajahafify/spritey/app/models"
)

type AssetsValidationResponse struct {
	OK       bool                            `json:"ok"`
	Assets   *models.AssetsValidationTarget  `json:"assets"`
	Pack     *models.Pack                    `json:"pack"`
	Summary  *models.AssetsValidationSummary `json:"summary"`
	Warnings []string                        `json:"warnings"`
	Errors   []models.Problem                `json:"errors"`
}

func WriteAssetsValidationJSON(writer io.Writer, result models.AssetsValidationResult) error {
	return writeJSON(writer, AssetsValidationResponse{
		OK:       true,
		Assets:   &result.Assets,
		Pack:     &result.Pack,
		Summary:  &result.Summary,
		Warnings: result.Warnings,
		Errors:   []models.Problem{},
	})
}

func WriteAssetsValidationErrorJSON(writer io.Writer, problem models.Problem) error {
	return writeJSON(writer, AssetsValidationResponse{
		OK:       false,
		Assets:   nil,
		Pack:     nil,
		Summary:  nil,
		Warnings: []string{},
		Errors:   []models.Problem{problem},
	})
}
