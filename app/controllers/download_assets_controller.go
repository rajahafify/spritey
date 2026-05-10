package controllers

import (
	"fmt"
	"io"
	"strings"

	"github.com/rajahafify/spritey/app/models"
	"github.com/rajahafify/spritey/app/services"
	"github.com/rajahafify/spritey/app/views"
)

type DownloadAssetsOptions struct {
	JSON  bool
	Force bool
}

type DownloadAssetsController struct {
	downloader services.LPCAssetsDownloader
}

var downloadAssetsProgressRendererFactory = newDownloadAssetsProgressRenderer

func NewDownloadAssetsController() DownloadAssetsController {
	return DownloadAssetsController{downloader: services.NewLPCAssetsDownloader()}
}

func (controller DownloadAssetsController) Download(options DownloadAssetsOptions, stdout io.Writer, stderr io.Writer) int {
	progress := downloadAssetsProgressRendererFactory(stderr)
	result, problem := controller.downloader.Download(services.DownloadLPCAssetsOptions{
		Force: options.Force,
		Progress: func(p services.DownloadProgress) {
			if options.JSON {
				return
			}
			progress(p)
		},
	})
	if problem != nil {
		if !options.JSON {
			fmt.Fprintln(stderr)
		}
		return WriteDownloadAssetsProblem(result.Assets.Path, *problem, downloadAssetsExitCode(problem.Code), options.JSON, stdout, stderr)
	}

	if !options.JSON {
		fmt.Fprintln(stderr)
	}
	if options.JSON {
		if err := views.WriteDownloadAssetsJSON(stdout, result); err != nil {
			fmt.Fprintln(stderr, err)
			return ExitGeneralError
		}
		return ExitSuccess
	}

	fmt.Fprintf(stdout, "ok: %s\n", result.Assets.Path)
	return ExitSuccess
}

func WriteDownloadAssetsProblem(assetsPath string, problem models.Problem, exitCode int, jsonMode bool, stdout io.Writer, stderr io.Writer) int {
	if jsonMode {
		if err := views.WriteDownloadAssetsErrorJSON(stdout, assetsPath, problem); err != nil {
			fmt.Fprintln(stderr, err)
			return ExitGeneralError
		}
		return exitCode
	}

	fmt.Fprintf(stderr, "%s: %s\n", problem.Code, problem.Message)
	return exitCode
}

func downloadAssetsExitCode(code string) int {
	switch code {
	case "ASSETS_DIRECTORY_NOT_FOUND", "READ_ASSETS_DIRECTORY_FAILED", "ASSETS_PATH_NOT_DIRECTORY",
		"MISSING_PACK_JSON", "READ_PACK_JSON_FAILED", "INVALID_PACK_JSON", "MISSING_SHEET_DEFINITIONS",
		"INVALID_SHEET_DEFINITION_JSON", "READ_SHEET_DEFINITIONS_FAILED", "MISSING_SPRITESHEETS", "MISSING_PALETTE_DEFINITIONS":
		return ExitInvalidAssets
	default:
		return ExitGeneralError
	}
}

func newDownloadAssetsProgressRenderer(writer io.Writer) func(services.DownloadProgress) {
	lastPercent := -1
	lastExtractCurrent := int64(-1)
	lastDownloadedMB := int64(-1)
	lastStage := ""

	writeLine := func(message string) {
		if _, err := fmt.Fprintf(writer, "\r%-78s", message); err == nil {
			lastStage = message
		}
	}

	return func(progress services.DownloadProgress) {
		switch progress.Stage {
		case "prepare":
			writeLine("Preparing LPC assets download...")
		case "download":
			if progress.Total > 0 {
				percent := int((progress.Current * 100) / progress.Total)
				if percent == lastPercent && !progress.Done {
					return
				}
				lastPercent = percent
				currentMB := float64(progress.Current) / (1024 * 1024)
				totalMB := float64(progress.Total) / (1024 * 1024)
				writeLine(fmt.Sprintf("Downloading LPC assets: %3d%% (%.1f/%.1f MB)", percent, currentMB, totalMB))
				return
			}
			currentMBFloor := progress.Current / (1024 * 1024)
			if currentMBFloor == lastDownloadedMB && !progress.Done {
				return
			}
			lastDownloadedMB = currentMBFloor
			writeLine(fmt.Sprintf("Downloading LPC assets: %.1f MB", float64(progress.Current)/(1024*1024)))
		case "extract":
			if progress.Total > 0 {
				if progress.Current == lastExtractCurrent && !progress.Done {
					return
				}
				lastExtractCurrent = progress.Current
				writeLine(fmt.Sprintf("Extracting assets: %d/%d", progress.Current, progress.Total))
				return
			}
			writeLine("Extracting assets...")
		case "validate":
			writeLine("Validating downloaded assets...")
		case "install":
			writeLine("Installing assets...")
		case "done":
			writeLine("LPC assets ready.")
		default:
			if strings.TrimSpace(lastStage) == "" {
				writeLine("Working...")
			}
		}
	}
}
