package controllers

import (
	"fmt"
	"io"

	"github.com/rajahafify/spritey/app/models"
	"github.com/rajahafify/spritey/app/services"
	"github.com/rajahafify/spritey/app/views"
)

const ExitValidationFailed = 5

type InspectLayerOptions struct {
	Target     string
	LayerID    string
	AssetsPath string
	JSON       bool
}

type InspectController struct {
	service services.InspectLayerService
}

func NewInspectController() InspectController {
	return InspectController{service: services.NewInspectLayerService()}
}

func (controller InspectController) InspectLayer(options InspectLayerOptions, stdout io.Writer, stderr io.Writer) int {
	if options.Target == "" {
		return WriteInspectProblem(
			models.Problem{
				Code:    "MISSING_INSPECT_TARGET",
				Message: "inspect target is required",
				Field:   "target",
			},
			ExitInvalidCLIUsage,
			options.JSON,
			stdout,
			stderr,
		)
	}
	if options.Target != "layer" {
		return WriteInspectProblem(
			models.Problem{
				Code:    "UNSUPPORTED_INSPECT_TARGET",
				Message: fmt.Sprintf("unsupported inspect target: %s", options.Target),
				Field:   "target",
			},
			ExitInvalidCLIUsage,
			options.JSON,
			stdout,
			stderr,
		)
	}
	if options.LayerID == "" {
		return WriteInspectProblem(
			models.Problem{
				Code:    "MISSING_LAYER_ID",
				Message: "layer id is required",
				Field:   "layer_id",
			},
			ExitInvalidCLIUsage,
			options.JSON,
			stdout,
			stderr,
		)
	}
	if options.AssetsPath == "" {
		return WriteInspectProblem(
			models.Problem{
				Code:    "MISSING_ASSETS",
				Message: "--assets is required",
				Field:   "assets",
			},
			ExitInvalidCLIUsage,
			options.JSON,
			stdout,
			stderr,
		)
	}

	result, problem := controller.service.Find(options.AssetsPath, options.LayerID)
	if problem != nil {
		exitCode := ExitInvalidAssets
		if problem.Code == "UNKNOWN_LAYER_ID" {
			exitCode = ExitValidationFailed
		}
		return WriteInspectProblem(*problem, exitCode, options.JSON, stdout, stderr)
	}

	if options.JSON {
		if err := views.WriteInspectLayerJSON(stdout, result); err != nil {
			fmt.Fprintln(stderr, err)
			return ExitGeneralError
		}
		return ExitSuccess
	}

	fmt.Fprintf(stdout, "%s/%s\n", result.Category, result.ID)
	return ExitSuccess
}

func WriteInspectProblem(problem models.Problem, exitCode int, jsonMode bool, stdout io.Writer, stderr io.Writer) int {
	if jsonMode {
		if err := views.WriteInspectLayerErrorJSON(stdout, problem); err != nil {
			fmt.Fprintln(stderr, err)
			return ExitGeneralError
		}
		return exitCode
	}

	fmt.Fprintf(stderr, "%s: %s\n", problem.Code, problem.Message)
	return exitCode
}
