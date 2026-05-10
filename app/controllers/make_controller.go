package controllers

import (
	"fmt"
	"io"

	"github.com/rajahafify/spritey/app/models"
	"github.com/rajahafify/spritey/app/services"
	"github.com/rajahafify/spritey/app/views"
)

const ExitRenderFailed = 6

type MakeOptions struct {
	RecipePath string
	AssetsPath string
	OutPath    string
	ReportPath string
	JSON       bool
}

type MakeController struct {
	service services.MakeService
}

func NewMakeController() MakeController {
	return MakeController{service: services.NewMakeService()}
}

func (controller MakeController) Make(options MakeOptions, stdout io.Writer, stderr io.Writer) int {
	if options.RecipePath == "" {
		return WriteMakeProblem(models.MakeOutputs{}, models.MakeProblem{
			Code:    "MISSING_RECIPE",
			Message: "recipe path is required",
			Field:   "recipe",
		}, ExitInvalidCLIUsage, options.JSON, stdout, stderr)
	}
	if options.AssetsPath == "" {
		return WriteMakeProblem(makeOutputs(options), models.MakeProblem{
			Code:    "MISSING_ASSETS",
			Message: "--assets is required",
			Field:   "assets",
		}, ExitInvalidCLIUsage, options.JSON, stdout, stderr)
	}
	if options.OutPath == "" {
		return WriteMakeProblem(makeOutputs(options), models.MakeProblem{
			Code:    "MISSING_OUT",
			Message: "--out is required",
			Field:   "out",
		}, ExitInvalidCLIUsage, options.JSON, stdout, stderr)
	}

	result, problem := controller.service.Make(options.RecipePath, options.AssetsPath, options.OutPath, options.ReportPath)
	if problem != nil {
		return WriteMakeProblem(makeOutputs(options), *problem, makeExitCode(problem.Code), options.JSON, stdout, stderr)
	}

	if options.JSON {
		if err := views.WriteMakeJSON(stdout, result); err != nil {
			fmt.Fprintln(stderr, err)
			return ExitGeneralError
		}
		return ExitSuccess
	}

	if err := views.WriteMakeText(stdout, result); err != nil {
		fmt.Fprintln(stderr, err)
		return ExitGeneralError
	}
	return ExitSuccess
}

func WriteMakeProblem(outputs models.MakeOutputs, problem models.MakeProblem, exitCode int, jsonMode bool, stdout io.Writer, stderr io.Writer) int {
	if jsonMode {
		if err := views.WriteMakeErrorJSON(stdout, outputs, problem); err != nil {
			fmt.Fprintln(stderr, err)
			return ExitGeneralError
		}
		return exitCode
	}
	fmt.Fprintf(stderr, "%s: %s\n", problem.Code, problem.Message)
	return exitCode
}

func makeExitCode(code string) int {
	switch code {
	case "RENDER_FAILED":
		return ExitRenderFailed
	case "RECIPE_FILE_NOT_FOUND", "READ_RECIPE_FAILED", "INVALID_RECIPE_JSON", "MISSING_SELECTIONS", "MISSING_SELECTION_ID",
		"UNKNOWN_LAYER_ID", "UNSUPPORTED_BODY_TYPE", "ASSETS_DIRECTORY_NOT_FOUND", "READ_ASSETS_DIRECTORY_FAILED", "ASSETS_PATH_NOT_DIRECTORY",
		"MISSING_PACK_JSON", "READ_PACK_JSON_FAILED", "INVALID_PACK_JSON", "MISSING_SHEET_DEFINITIONS", "INVALID_SHEET_DEFINITION_JSON", "READ_SHEET_DEFINITIONS_FAILED":
		return ExitInvalidAssets
	default:
		return ExitGeneralError
	}
}

func makeOutputs(options MakeOptions) models.MakeOutputs {
	outputs := models.MakeOutputs{
		PNG: models.MakeFileOutput{Path: options.OutPath},
	}
	if options.ReportPath != "" {
		outputs.Report = &models.MakeFileOutput{Path: options.ReportPath}
	}
	return outputs
}
