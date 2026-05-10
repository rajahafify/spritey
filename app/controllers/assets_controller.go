package controllers

import (
	"fmt"
	"io"

	"github.com/rajahafify/spritey/app/models"
	"github.com/rajahafify/spritey/app/services"
	"github.com/rajahafify/spritey/app/views"
)

type AssetsValidateOptions struct {
	Subcommand string
	AssetsPath string
	JSON       bool
}

type AssetsController struct {
	validator services.AssetsValidator
}

func NewAssetsController() AssetsController {
	return AssetsController{validator: services.NewAssetsValidator()}
}

func (controller AssetsController) Validate(options AssetsValidateOptions, stdout io.Writer, stderr io.Writer) int {
	if options.Subcommand == "" {
		return WriteAssetsValidationProblem(
			models.Problem{
				Code:    "MISSING_ASSETS_SUBCOMMAND",
				Message: "assets subcommand is required",
				Field:   "subcommand",
			},
			ExitInvalidCLIUsage,
			options.JSON,
			stdout,
			stderr,
		)
	}
	if options.Subcommand != "validate" {
		return WriteAssetsValidationProblem(
			models.Problem{
				Code:    "UNSUPPORTED_ASSETS_SUBCOMMAND",
				Message: fmt.Sprintf("unsupported assets subcommand: %s", options.Subcommand),
				Field:   "subcommand",
			},
			ExitInvalidCLIUsage,
			options.JSON,
			stdout,
			stderr,
		)
	}
	if options.AssetsPath == "" {
		return WriteAssetsValidationProblem(
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

	result, problem := controller.validator.Validate(options.AssetsPath)
	if problem != nil {
		return WriteAssetsValidationProblem(*problem, ExitInvalidAssets, options.JSON, stdout, stderr)
	}

	if options.JSON {
		if err := views.WriteAssetsValidationJSON(stdout, result); err != nil {
			fmt.Fprintln(stderr, err)
			return ExitGeneralError
		}
		return ExitSuccess
	}

	fmt.Fprintf(stdout, "ok: %s\n", result.Assets.Path)
	return ExitSuccess
}

func WriteAssetsValidationProblem(problem models.Problem, exitCode int, jsonMode bool, stdout io.Writer, stderr io.Writer) int {
	if jsonMode {
		if err := views.WriteAssetsValidationErrorJSON(stdout, problem); err != nil {
			fmt.Fprintln(stderr, err)
			return ExitGeneralError
		}
		return exitCode
	}

	fmt.Fprintf(stderr, "%s: %s\n", problem.Code, problem.Message)
	return exitCode
}
