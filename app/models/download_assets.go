package models

type DownloadAssetsTarget struct {
	Path    string `json:"path"`
	Source  string `json:"source,omitempty"`
	Version string `json:"version,omitempty"`
}

type DownloadAssetsResult struct {
	Assets   DownloadAssetsTarget `json:"assets"`
	Warnings []string             `json:"warnings"`
}
