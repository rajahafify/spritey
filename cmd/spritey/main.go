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
	case "assets":
		options, problem := parseAssetsValidateOptions(args[1:])
		if problem != nil {
			return controllers.WriteAssetsValidationProblem(*problem, controllers.ExitInvalidCLIUsage, options.JSON, stdout, stderr)
		}
		return controllers.NewAssetsController().Validate(options, stdout, stderr)
	case "catalog":
		options, problem := parseCatalogOptions(args[1:])
		if problem != nil {
			return controllers.WriteProblem(*problem, controllers.ExitInvalidCLIUsage, options.JSON, stdout, stderr)
		}
		return controllers.NewCatalogController().Catalog(options, stdout, stderr)
	case "inspect":
		options, problem := parseInspectOptions(args[1:])
		if problem != nil {
			return controllers.WriteInspectProblem(*problem, controllers.ExitInvalidCLIUsage, options.JSON, stdout, stderr)
		}
		return controllers.NewInspectController().InspectLayer(options, stdout, stderr)
	case "validate":
		options, problem := parseValidateOptions(args[1:])
		if problem != nil {
			return controllers.WriteValidateProblem(*problem, controllers.ExitInvalidCLIUsage, options.JSON, stdout, stderr)
		}
		return controllers.NewValidateController().Validate(options, stdout, stderr)
	case "make":
		options, problem := parseMakeOptions(args[1:])
		if problem != nil {
			return controllers.WriteMakeProblem(models.MakeOutputs{}, models.MakeProblem{
				Code:    problem.Code,
				Message: problem.Message,
				Field:   problem.Field,
			}, controllers.ExitInvalidCLIUsage, options.JSON, stdout, stderr)
		}
		return controllers.NewMakeController().Make(options, stdout, stderr)
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

func parseMakeOptions(args []string) (controllers.MakeOptions, *models.Problem) {
	options := controllers.MakeOptions{}
	if len(args) == 0 {
		return options, nil
	}
	if !isFlag(args[0]) {
		options.RecipePath = args[0]
		args = args[1:]
	}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
			options.JSON = true
		case "--assets":
			if i+1 >= len(args) {
				return options, &models.Problem{Code: "MISSING_ASSETS_VALUE", Message: "--assets requires a value", Field: "assets"}
			}
			options.AssetsPath = args[i+1]
			i++
		case "--out":
			if i+1 >= len(args) {
				return options, &models.Problem{Code: "MISSING_OUT_VALUE", Message: "--out requires a value", Field: "out"}
			}
			options.OutPath = args[i+1]
			i++
		case "--report":
			if i+1 >= len(args) {
				return options, &models.Problem{Code: "MISSING_REPORT_VALUE", Message: "--report requires a value", Field: "report"}
			}
			options.ReportPath = args[i+1]
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

func parseAssetsValidateOptions(args []string) (controllers.AssetsValidateOptions, *models.Problem) {
	options := controllers.AssetsValidateOptions{}
	if len(args) == 0 || isFlag(args[0]) {
		options.JSON = hasJSONFlag(args)
		return options, nil
	}

	options.Subcommand = args[0]
	return parseAssetsValidateFlags(options, args[1:])
}

func parseAssetsValidateFlags(options controllers.AssetsValidateOptions, args []string) (controllers.AssetsValidateOptions, *models.Problem) {
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

func parseValidateOptions(args []string) (controllers.ValidateOptions, *models.Problem) {
	options := controllers.ValidateOptions{}
	if len(args) == 0 {
		return options, nil
	}
	if isFlag(args[0]) {
		return parseValidateFlags(options, args)
	}
	options.RecipePath = args[0]
	return parseValidateFlags(options, args[1:])
}

func parseValidateFlags(options controllers.ValidateOptions, args []string) (controllers.ValidateOptions, *models.Problem) {
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

func parseInspectOptions(args []string) (controllers.InspectLayerOptions, *models.Problem) {
	options := controllers.InspectLayerOptions{}
	if len(args) == 0 {
		return options, nil
	}

	if isFlag(args[0]) {
		return parseInspectFlags(options, args)
	}
	options.Target = args[0]
	args = args[1:]

	if options.Target == "layer" {
		if len(args) == 0 || isFlag(args[0]) {
			return parseInspectFlags(options, args)
		}
		options.LayerID = args[0]
		args = args[1:]
	}

	return parseInspectFlags(options, args)
}

func parseInspectFlags(options controllers.InspectLayerOptions, args []string) (controllers.InspectLayerOptions, *models.Problem) {
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

func isFlag(arg string) bool {
	return len(arg) > 0 && arg[0] == '-'
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
