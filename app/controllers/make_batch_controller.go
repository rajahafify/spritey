package controllers

import (
	"fmt"
	"io"

	"github.com/rajahafify/spritey/app/models"
	"github.com/rajahafify/spritey/app/services"
	"github.com/rajahafify/spritey/app/views"
)

type MakeBatchOptions struct {
	ManifestPath string
	AssetsPath   string
	JSON         bool
}

type MakeBatchController struct {
	service services.MakeBatchService
}

func NewMakeBatchController() MakeBatchController {
	return MakeBatchController{service: services.NewMakeBatchService()}
}

func (controller MakeBatchController) Make(options MakeBatchOptions, stdout io.Writer, stderr io.Writer) int {
	if options.ManifestPath == "" {
		return WriteMakeBatchProblem(models.MakeBatchResult{}, models.MakeProblem{
			Code:    "MISSING_MANIFEST",
			Message: "manifest path is required",
			Field:   "manifest",
		}, ExitInvalidCLIUsage, options.JSON, stdout, stderr)
	}
	if options.AssetsPath == "" {
		return WriteMakeBatchProblem(models.MakeBatchResult{}, models.MakeProblem{
			Code:    "MISSING_ASSETS",
			Message: "--assets is required",
			Field:   "assets",
		}, ExitInvalidCLIUsage, options.JSON, stdout, stderr)
	}

	result, problem := controller.service.Make(options.ManifestPath, options.AssetsPath)
	if problem != nil {
		return WriteMakeBatchProblem(result, *problem, makeBatchExitCode(problem.Code), options.JSON, stdout, stderr)
	}

	if options.JSON {
		if err := views.WriteMakeBatchJSON(stdout, result); err != nil {
			fmt.Fprintln(stderr, err)
			return ExitGeneralError
		}
		return ExitSuccess
	}

	if err := views.WriteMakeBatchText(stdout, result); err != nil {
		fmt.Fprintln(stderr, err)
		return ExitGeneralError
	}
	return ExitSuccess
}

func WriteMakeBatchProblem(result models.MakeBatchResult, problem models.MakeProblem, exitCode int, jsonMode bool, stdout io.Writer, stderr io.Writer) int {
	if jsonMode {
		if err := views.WriteMakeBatchErrorJSON(stdout, result, problem); err != nil {
			fmt.Fprintln(stderr, err)
			return ExitGeneralError
		}
		return exitCode
	}
	fmt.Fprintf(stderr, "%s: %s\n", problem.Code, problem.Message)
	return exitCode
}

func makeBatchExitCode(code string) int {
	switch code {
	case "BATCH_MANIFEST_FILE_NOT_FOUND", "READ_BATCH_MANIFEST_FAILED", "INVALID_BATCH_MANIFEST_JSON",
		"UNSUPPORTED_BATCH_MANIFEST_SCHEMA", "EMPTY_BATCH_JOBS", "MISSING_BATCH_JOB_ID",
		"MISSING_BATCH_JOB_RECIPE", "MISSING_BATCH_JOB_OUT":
		return ExitInvalidRecipe
	default:
		return makeCommandExitCode(code)
	}
}
