package controllers

import (
	"fmt"
	"io"

	"github.com/rajahafify/spritey/app/models"
	"github.com/rajahafify/spritey/app/services"
	"github.com/rajahafify/spritey/app/views"
)

const ExitInvalidRecipe = 4

type ValidateOptions struct {
	RecipePath string
	AssetsPath string
	JSON       bool
}

type ValidateController struct {
	validator services.RecipeValidator
}

func NewValidateController() ValidateController {
	return ValidateController{validator: services.NewRecipeValidator()}
}

func (controller ValidateController) Validate(options ValidateOptions, stdout io.Writer, stderr io.Writer) int {
	if options.RecipePath == "" {
		return WriteValidateProblem(
			models.Problem{
				Code:    "MISSING_RECIPE",
				Message: "recipe path is required",
				Field:   "recipe",
			},
			ExitInvalidCLIUsage,
			options.JSON,
			stdout,
			stderr,
		)
	}
	if options.AssetsPath == "" {
		return WriteValidateProblem(
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

	result, problem := controller.validator.Validate(options.RecipePath, options.AssetsPath)
	if problem != nil {
		return WriteValidateProblem(*problem, validateExitCode(problem.Code), options.JSON, stdout, stderr)
	}

	if options.JSON {
		if err := views.WriteValidateRecipeJSON(stdout, result); err != nil {
			fmt.Fprintln(stderr, err)
			return ExitGeneralError
		}
		return ExitSuccess
	}

	fmt.Fprintf(stdout, "ok: %s\n", result.Path)
	return ExitSuccess
}

func validateExitCode(code string) int {
	switch code {
	case "RECIPE_FILE_NOT_FOUND", "READ_RECIPE_FAILED", "INVALID_RECIPE_JSON", "MISSING_SELECTIONS", "MISSING_SELECTION_ID":
		return ExitInvalidRecipe
	case "UNKNOWN_LAYER_ID", "UNSUPPORTED_BODY_TYPE":
		return ExitValidationFailed
	case "ASSETS_DIRECTORY_NOT_FOUND", "READ_ASSETS_DIRECTORY_FAILED", "ASSETS_PATH_NOT_DIRECTORY", "MISSING_PACK_JSON", "READ_PACK_JSON_FAILED", "INVALID_PACK_JSON", "MISSING_SHEET_DEFINITIONS", "INVALID_SHEET_DEFINITION_JSON", "READ_SHEET_DEFINITIONS_FAILED":
		return ExitInvalidAssets
	default:
		return ExitGeneralError
	}
}

func WriteValidateProblem(problem models.Problem, exitCode int, jsonMode bool, stdout io.Writer, stderr io.Writer) int {
	if jsonMode {
		if err := views.WriteValidateRecipeErrorJSON(stdout, problem); err != nil {
			fmt.Fprintln(stderr, err)
			return ExitGeneralError
		}
		return exitCode
	}

	fmt.Fprintf(stderr, "%s: %s\n", problem.Code, problem.Message)
	return exitCode
}
