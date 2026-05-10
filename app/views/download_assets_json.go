package views

import (
	"io"

	"github.com/rajahafify/spritey/app/models"
)

type DownloadAssetsResponse struct {
	OK       bool                         `json:"ok"`
	Command  string                       `json:"command"`
	Assets   *models.DownloadAssetsTarget `json:"assets"`
	Warnings []string                     `json:"warnings"`
	Errors   []models.Problem             `json:"errors"`
}

func WriteDownloadAssetsJSON(writer io.Writer, result models.DownloadAssetsResult) error {
	return writeJSON(writer, DownloadAssetsResponse{
		OK:       true,
		Command:  "download-lpc-assets",
		Assets:   &result.Assets,
		Warnings: result.Warnings,
		Errors:   []models.Problem{},
	})
}

func WriteDownloadAssetsErrorJSON(writer io.Writer, assetsPath string, problem models.Problem) error {
	var assets *models.DownloadAssetsTarget
	if assetsPath != "" {
		assets = &models.DownloadAssetsTarget{Path: assetsPath}
	}
	return writeJSON(writer, DownloadAssetsResponse{
		OK:       false,
		Command:  "download-lpc-assets",
		Assets:   assets,
		Warnings: []string{},
		Errors:   []models.Problem{problem},
	})
}
