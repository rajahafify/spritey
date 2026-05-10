package controllers

import (
	"fmt"
	"io"

	"github.com/rajahafify/spritey/app/models"
	"github.com/rajahafify/spritey/app/services"
	"github.com/rajahafify/spritey/app/views"
)

const (
	ExitSuccess         = 0
	ExitGeneralError    = 1
	ExitInvalidCLIUsage = 2
	ExitInvalidAssets   = 3
)

type CatalogOptions struct {
	AssetsPath string
	JSON       bool
}

type CatalogController struct {
	loader services.CatalogLoader
}

func NewCatalogController() CatalogController {
	return CatalogController{loader: services.NewCatalogLoader()}
}

func (controller CatalogController) Catalog(options CatalogOptions, stdout io.Writer, stderr io.Writer) int {
	if options.AssetsPath == "" {
		return WriteProblem(
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

	catalog, loadErr := controller.loader.Load(options.AssetsPath)
	if loadErr != nil {
		return WriteProblem(*loadErr, ExitInvalidAssets, options.JSON, stdout, stderr)
	}

	if options.JSON {
		if err := views.WriteCatalogJSON(stdout, catalog); err != nil {
			fmt.Fprintln(stderr, err)
			return ExitGeneralError
		}
		return ExitSuccess
	}

	for _, category := range catalog.Categories {
		fmt.Fprintf(stdout, "%s:", category.ID)
		for _, layer := range category.Layers {
			fmt.Fprintf(stdout, " %s", layer.ID)
		}
		fmt.Fprintln(stdout)
	}
	return ExitSuccess
}

func WriteProblem(problem models.Problem, exitCode int, jsonMode bool, stdout io.Writer, stderr io.Writer) int {
	if jsonMode {
		if err := views.WriteCatalogErrorJSON(stdout, problem); err != nil {
			fmt.Fprintln(stderr, err)
			return ExitGeneralError
		}
		return exitCode
	}

	fmt.Fprintf(stderr, "%s: %s\n", problem.Code, problem.Message)
	return exitCode
}
