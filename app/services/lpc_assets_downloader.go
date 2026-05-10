package services

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rajahafify/spritey/app/models"
)

const (
	lpcAssetsSourceRepository = "https://github.com/liberatedpixelcup/Universal-LPC-Spritesheet-Character-Generator"
	lpcAssetsArchiveURL       = "https://codeload.github.com/liberatedpixelcup/Universal-LPC-Spritesheet-Character-Generator/zip/refs/heads/master"
	lpcAssetsVersion          = "master"
)

type DownloadProgress struct {
	Stage   string
	Current int64
	Total   int64
	Done    bool
}

type DownloadProgressCallback func(DownloadProgress)
type LPCArchiveFetcher func(url string, onProgress func(received int64, total int64)) ([]byte, error)
type UserCacheDirResolver func() (string, error)

var (
	DefaultLPCArchiveFetcher    LPCArchiveFetcher    = fetchLPCArchive
	DefaultUserCacheDirResolver UserCacheDirResolver = os.UserCacheDir
)

type DownloadLPCAssetsOptions struct {
	Force    bool
	Progress DownloadProgressCallback
}

type LPCAssetsDownloader struct {
	archiveFetcher       LPCArchiveFetcher
	userCacheDirResolver UserCacheDirResolver
	validator            AssetsValidator
	source               string
	archiveURL           string
	version              string
}

func NewLPCAssetsDownloader() LPCAssetsDownloader {
	return LPCAssetsDownloader{
		archiveFetcher:       DefaultLPCArchiveFetcher,
		userCacheDirResolver: DefaultUserCacheDirResolver,
		validator:            NewAssetsValidator(),
		source:               lpcAssetsSourceRepository,
		archiveURL:           lpcAssetsArchiveURL,
		version:              lpcAssetsVersion,
	}
}

func NewLPCAssetsDownloaderWithDeps(fetcher LPCArchiveFetcher, cacheResolver UserCacheDirResolver) LPCAssetsDownloader {
	downloader := NewLPCAssetsDownloader()
	if fetcher != nil {
		downloader.archiveFetcher = fetcher
	}
	if cacheResolver != nil {
		downloader.userCacheDirResolver = cacheResolver
	}
	return downloader
}

func (downloader LPCAssetsDownloader) Download(options DownloadLPCAssetsOptions) (models.DownloadAssetsResult, *models.Problem) {
	emitProgress := func(stage string, current int64, total int64, done bool) {
		if options.Progress != nil {
			options.Progress(DownloadProgress{
				Stage:   stage,
				Current: current,
				Total:   total,
				Done:    done,
			})
		}
	}

	emitProgress("prepare", 0, 0, false)

	cacheDir, err := downloader.userCacheDirResolver()
	if err != nil {
		return models.DownloadAssetsResult{}, &models.Problem{
			Code:    "RESOLVE_CACHE_DIR_FAILED",
			Message: err.Error(),
			Field:   "assets",
		}
	}

	targetPath := filepath.Join(cacheDir, "spritey", "assets", "lpc")
	result := models.DownloadAssetsResult{
		Assets: models.DownloadAssetsTarget{
			Path:    targetPath,
			Source:  downloader.source,
			Version: downloader.version,
		},
		Warnings: []string{},
	}

	if !options.Force {
		if existing, statErr := os.Stat(targetPath); statErr == nil && existing.IsDir() {
			if _, validateErr := downloader.validator.Validate(targetPath); validateErr == nil {
				emitProgress("done", 1, 1, true)
				return result, nil
			}
		}
	}

	emitProgress("download", 0, 0, false)
	archiveBytes, err := downloader.archiveFetcher(downloader.archiveURL, func(received int64, total int64) {
		emitProgress("download", received, total, false)
	})
	if err != nil {
		return result, &models.Problem{
			Code:    "DOWNLOAD_FAILED",
			Message: err.Error(),
			Field:   "assets",
		}
	}
	emitProgress("download", int64(len(archiveBytes)), int64(len(archiveBytes)), true)

	stagePath := targetPath + ".staging"
	if err := os.RemoveAll(stagePath); err != nil {
		return result, &models.Problem{
			Code:    "WRITE_ASSETS_FAILED",
			Message: err.Error(),
			Field:   "assets",
		}
	}
	if err := os.MkdirAll(stagePath, 0o755); err != nil {
		return result, &models.Problem{
			Code:    "WRITE_ASSETS_FAILED",
			Message: err.Error(),
			Field:   "assets",
		}
	}
	defer os.RemoveAll(stagePath)

	emitProgress("extract", 0, 0, false)
	foundPack, problem := extractLPCArchiveSubset(archiveBytes, stagePath, func(current int64, total int64) {
		emitProgress("extract", current, total, current == total && total > 0)
	})
	if problem != nil {
		return result, problem
	}
	if !foundPack {
		if problem := writeDefaultPackJSON(filepath.Join(stagePath, "pack.json")); problem != nil {
			return result, problem
		}
	}
	for _, requiredDir := range []string{"sheet_definitions", "spritesheets", "palette_definitions"} {
		if err := os.MkdirAll(filepath.Join(stagePath, requiredDir), 0o755); err != nil {
			return result, &models.Problem{
				Code:    "WRITE_ASSETS_FAILED",
				Message: err.Error(),
				Field:   "assets",
			}
		}
	}

	emitProgress("validate", 0, 0, false)
	if _, validateErr := downloader.validator.Validate(stagePath); validateErr != nil {
		return result, validateErr
	}
	emitProgress("validate", 1, 1, true)

	if err := os.RemoveAll(targetPath); err != nil {
		return result, &models.Problem{
			Code:    "WRITE_ASSETS_FAILED",
			Message: err.Error(),
			Field:   "assets",
		}
	}
	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return result, &models.Problem{
			Code:    "WRITE_ASSETS_FAILED",
			Message: err.Error(),
			Field:   "assets",
		}
	}
	if err := os.Rename(stagePath, targetPath); err != nil {
		return result, &models.Problem{
			Code:    "WRITE_ASSETS_FAILED",
			Message: err.Error(),
			Field:   "assets",
		}
	}

	emitProgress("install", 1, 1, true)
	emitProgress("done", 1, 1, true)
	return result, nil
}

func extractLPCArchiveSubset(archiveBytes []byte, destination string, onProgress func(current int64, total int64)) (bool, *models.Problem) {
	reader, err := zip.NewReader(bytes.NewReader(archiveBytes), int64(len(archiveBytes)))
	if err != nil {
		return false, &models.Problem{
			Code:    "EXTRACT_FAILED",
			Message: err.Error(),
			Field:   "assets",
		}
	}

	foundPack := false
	filesToExtract := make([]*zip.File, 0, len(reader.File))
	for _, file := range reader.File {
		relativePath := trimArchiveRoot(file.Name)
		if relativePath == "" {
			continue
		}
		relativePath = filepath.ToSlash(relativePath)
		if !shouldExtractArchivePath(relativePath) {
			continue
		}
		filesToExtract = append(filesToExtract, file)
	}

	totalFiles := int64(len(filesToExtract))
	if onProgress != nil && totalFiles > 0 {
		onProgress(0, totalFiles)
	}

	for i, file := range filesToExtract {
		relativePath := filepath.ToSlash(trimArchiveRoot(file.Name))
		if relativePath == "pack.json" {
			foundPack = true
		}

		targetFilePath := filepath.Join(destination, filepath.FromSlash(relativePath))
		if !isWithinBasePath(destination, targetFilePath) {
			return false, &models.Problem{
				Code:    "EXTRACT_FAILED",
				Message: fmt.Sprintf("invalid archive path: %s", file.Name),
				Field:   "assets",
			}
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetFilePath, 0o755); err != nil {
				return false, &models.Problem{
					Code:    "WRITE_ASSETS_FAILED",
					Message: err.Error(),
					Field:   "assets",
				}
			}
			if onProgress != nil {
				onProgress(int64(i+1), totalFiles)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetFilePath), 0o755); err != nil {
			return false, &models.Problem{
				Code:    "WRITE_ASSETS_FAILED",
				Message: err.Error(),
				Field:   "assets",
			}
		}

		sourceFile, err := file.Open()
		if err != nil {
			return false, &models.Problem{
				Code:    "EXTRACT_FAILED",
				Message: err.Error(),
				Field:   "assets",
			}
		}
		targetFile, err := os.Create(targetFilePath)
		if err != nil {
			sourceFile.Close()
			return false, &models.Problem{
				Code:    "WRITE_ASSETS_FAILED",
				Message: err.Error(),
				Field:   "assets",
			}
		}
		if _, err := io.Copy(targetFile, sourceFile); err != nil {
			targetFile.Close()
			sourceFile.Close()
			return false, &models.Problem{
				Code:    "WRITE_ASSETS_FAILED",
				Message: err.Error(),
				Field:   "assets",
			}
		}
		targetFile.Close()
		sourceFile.Close()

		if onProgress != nil {
			onProgress(int64(i+1), totalFiles)
		}
	}

	return foundPack, nil
}

func shouldExtractArchivePath(path string) bool {
	return path == "pack.json" ||
		strings.HasPrefix(path, "sheet_definitions/") ||
		strings.HasPrefix(path, "spritesheets/") ||
		strings.HasPrefix(path, "palette_definitions/")
}

func trimArchiveRoot(path string) string {
	parts := strings.Split(filepath.ToSlash(path), "/")
	if len(parts) <= 1 {
		return ""
	}
	return strings.Join(parts[1:], "/")
}

func isWithinBasePath(base string, target string) bool {
	baseClean := filepath.Clean(base)
	targetClean := filepath.Clean(target)
	if baseClean == targetClean {
		return true
	}
	return strings.HasPrefix(targetClean, baseClean+string(os.PathSeparator))
}

func writeDefaultPackJSON(path string) *models.Problem {
	defaultPack := models.Pack{
		SchemaVersion: "1",
		ID:            "lpc-assets",
		Name:          "LPC Assets",
		Defaults: models.PackDefaults{
			BodyType:                "male",
			Animations:              []string{"spellcast", "thrust", "walk", "slash", "shoot", "hurt"},
			CanvasWidth:             832,
			MissingBodyTypeFallback: "male",
			PaletteSourceFallbacks:  []string{"light", "base", "beige", "brown"},
		},
	}

	data, err := json.MarshalIndent(defaultPack, "", "  ")
	if err != nil {
		return &models.Problem{
			Code:    "WRITE_ASSETS_FAILED",
			Message: err.Error(),
			Field:   "assets",
		}
	}
	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return &models.Problem{
			Code:    "WRITE_ASSETS_FAILED",
			Message: err.Error(),
			Field:   "assets",
		}
	}
	return nil
}

func fetchLPCArchive(url string, onProgress func(received int64, total int64)) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", response.Status)
	}

	total := response.ContentLength
	if onProgress != nil {
		onProgress(0, total)
	}

	var output bytes.Buffer
	buffer := make([]byte, 64*1024)
	var received int64
	for {
		n, readErr := response.Body.Read(buffer)
		if n > 0 {
			if _, err := output.Write(buffer[:n]); err != nil {
				return nil, err
			}
			received += int64(n)
			if onProgress != nil {
				onProgress(received, total)
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return nil, readErr
		}
	}

	return output.Bytes(), nil
}
