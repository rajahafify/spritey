package views

import (
	"io"

	"github.com/rajahafify/spritey/app/models"
)

type MakeBatchResponse struct {
	OK       bool                        `json:"ok"`
	Command  string                      `json:"command"`
	Summary  models.MakeBatchSummary     `json:"summary"`
	Jobs     []models.MakeBatchJobResult `json:"jobs"`
	Warnings []string                    `json:"warnings"`
	Errors   []models.MakeProblem        `json:"errors"`
}

func WriteMakeBatchJSON(writer io.Writer, result models.MakeBatchResult) error {
	return writeJSON(writer, MakeBatchResponse{
		OK:       true,
		Command:  "make-batch",
		Summary:  result.Summary,
		Jobs:     result.Jobs,
		Warnings: result.Warnings,
		Errors:   []models.MakeProblem{},
	})
}

func WriteMakeBatchErrorJSON(writer io.Writer, result models.MakeBatchResult, problem models.MakeProblem) error {
	return writeJSON(writer, MakeBatchResponse{
		OK:       false,
		Command:  "make-batch",
		Summary:  result.Summary,
		Jobs:     result.Jobs,
		Warnings: result.Warnings,
		Errors:   []models.MakeProblem{problem},
	})
}
