package main

import (
	"fmt"
	"io"
	"os"

	"github.com/rajahafify/spritey/app/controllers"
	"github.com/rajahafify/spritey/app/models"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		return controllers.WriteProblem(
			models.Problem{
				Code:    "MISSING_COMMAND",
				Message: "command is required",
				Field:   "command",
			},
			controllers.ExitInvalidCLIUsage,
			false,
			stdout,
			stderr,
		)
	}

	switch args[0] {
	case "catalog":
		options, problem := parseCatalogOptions(args[1:])
		if problem != nil {
			return controllers.WriteProblem(*problem, controllers.ExitInvalidCLIUsage, options.JSON, stdout, stderr)
		}
		return controllers.NewCatalogController().Catalog(options, stdout, stderr)
	default:
		return controllers.WriteProblem(
			models.Problem{
				Code:    "UNKNOWN_COMMAND",
				Message: fmt.Sprintf("unknown command: %s", args[0]),
				Field:   "command",
			},
			controllers.ExitInvalidCLIUsage,
			hasJSONFlag(args),
			stdout,
			stderr,
		)
	}
}

func parseCatalogOptions(args []string) (controllers.CatalogOptions, *models.Problem) {
	options := controllers.CatalogOptions{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
			options.JSON = true
		case "--assets":
			if i+1 >= len(args) {
				return options, &models.Problem{
					Code:    "MISSING_ASSETS_VALUE",
					Message: "--assets requires a value",
					Field:   "assets",
				}
			}
			options.AssetsPath = args[i+1]
			i++
		default:
			return options, &models.Problem{
				Code:    "UNKNOWN_ARGUMENT",
				Message: fmt.Sprintf("unknown argument: %s", args[i]),
				Field:   "argument",
			}
		}
	}
	return options, nil
}

func hasJSONFlag(args []string) bool {
	for _, arg := range args {
		if arg == "--json" {
			return true
		}
	}
	return false
}
