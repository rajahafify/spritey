package views

import (
	"io"

	"github.com/rajahafify/spritey/app/models"
)

type MakeResponse struct {
	OK       bool                   `json:"ok"`
	Command  string                 `json:"command"`
	Outputs  models.MakeOutputs     `json:"outputs"`
	Summary  map[string]interface{} `json:"summary"`
	Warnings []string               `json:"warnings"`
	Errors   []models.MakeProblem   `json:"errors"`
}

func WriteMakeJSON(writer io.Writer, result models.MakeResult) error {
	return writeJSON(writer, MakeResponse{
		OK:      true,
		Command: "make",
		Outputs: result.Outputs,
		Summary: map[string]interface{}{
			"frame_count": result.Summary.FrameCount,
			"canvas": map[string]interface{}{
				"width":  result.Summary.Canvas.Width,
				"height": result.Summary.Canvas.Height,
			},
			"animation_count": result.Summary.AnimationCount,
		},
		Warnings: result.Warnings,
		Errors:   []models.MakeProblem{},
	})
}

func WriteMakeErrorJSON(writer io.Writer, outputs models.MakeOutputs, problem models.MakeProblem) error {
	return writeJSON(writer, MakeResponse{
		OK:       false,
		Command:  "make",
		Outputs:  outputs,
		Summary:  map[string]interface{}{},
		Warnings: []string{},
		Errors:   []models.MakeProblem{problem},
	})
}
