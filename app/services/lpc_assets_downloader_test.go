package services

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLPCAssetsDownloaderUsesUserCacheDirPath(t *testing.T) {
	cacheRoot := t.TempDir()
	archivePayload := makeLPCArchivePayload(t, false, "v1")

	downloader := NewLPCAssetsDownloaderWithDeps(
		func(string, func(int64, int64)) ([]byte, error) { return archivePayload, nil },
		func() (string, error) { return cacheRoot, nil },
	)

	result, problem := downloader.Download(DownloadLPCAssetsOptions{})
	if problem != nil {
		t.Fatalf("expected success, got problem %+v", problem)
	}

	expectedPath := filepath.Join(cacheRoot, "spritey", "assets", "lpc")
	if result.Assets.Path != expectedPath {
		t.Fatalf("expected assets path %q, got %q", expectedPath, result.Assets.Path)
	}

	if _, err := os.Stat(filepath.Join(expectedPath, "pack.json")); err != nil {
		t.Fatalf("expected generated pack.json, stat error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(expectedPath, "sheet_definitions", "body", "body_human.json")); err != nil {
		t.Fatalf("expected sheet definition, stat error: %v", err)
	}

	if _, validateErr := NewAssetsValidator().Validate(expectedPath); validateErr != nil {
		t.Fatalf("expected installed assets to validate, got %+v", validateErr)
	}
}

func TestLPCAssetsDownloaderInvalidDownloadedAssets(t *testing.T) {
	cacheRoot := t.TempDir()
	archivePayload := makeLPCArchivePayload(t, true, "v1")

	downloader := NewLPCAssetsDownloaderWithDeps(
		func(string, func(int64, int64)) ([]byte, error) { return archivePayload, nil },
		func() (string, error) { return cacheRoot, nil },
	)

	_, problem := downloader.Download(DownloadLPCAssetsOptions{})
	if problem == nil {
		t.Fatal("expected validation problem")
	}
	if problem.Code != "INVALID_PACK_JSON" {
		t.Fatalf("expected INVALID_PACK_JSON, got %+v", problem)
	}
}

func TestLPCAssetsDownloaderForceReinstalls(t *testing.T) {
	cacheRoot := t.TempDir()
	fetchCount := 0

	downloader := NewLPCAssetsDownloaderWithDeps(
		func(string, func(int64, int64)) ([]byte, error) {
			fetchCount++
			return makeLPCArchivePayload(t, false, "v"+string(rune('0'+fetchCount))), nil
		},
		func() (string, error) { return cacheRoot, nil },
	)

	firstResult, problem := downloader.Download(DownloadLPCAssetsOptions{})
	if problem != nil {
		t.Fatalf("expected first install success, got %+v", problem)
	}
	if fetchCount != 1 {
		t.Fatalf("expected one fetch on first download, got %d", fetchCount)
	}

	secondResult, problem := downloader.Download(DownloadLPCAssetsOptions{})
	if problem != nil {
		t.Fatalf("expected idempotent second install success, got %+v", problem)
	}
	if secondResult.Assets.Path != firstResult.Assets.Path {
		t.Fatalf("expected same target path, got %q vs %q", secondResult.Assets.Path, firstResult.Assets.Path)
	}
	if fetchCount != 1 {
		t.Fatalf("expected second download without force to skip fetch, got %d fetches", fetchCount)
	}

	_, problem = downloader.Download(DownloadLPCAssetsOptions{Force: true})
	if problem != nil {
		t.Fatalf("expected force reinstall success, got %+v", problem)
	}
	if fetchCount != 2 {
		t.Fatalf("expected force reinstall to refetch archive, got %d fetches", fetchCount)
	}
}

func makeLPCArchivePayload(t *testing.T, includeInvalidPack bool, marker string) []byte {
	t.Helper()
	buffer := &bytes.Buffer{}
	archive := zip.NewWriter(buffer)

	base := "Universal-LPC-Spritesheet-Character-Generator-master/"
	addFile := func(path string, body string) {
		writer, err := archive.Create(base + strings.TrimPrefix(filepath.ToSlash(path), "/"))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.WriteString(writer, body); err != nil {
			t.Fatal(err)
		}
	}

	if includeInvalidPack {
		addFile("pack.json", `{not-json}`)
	}
	addFile("sheet_definitions/body/body_human.json", `{"name":"Human Body","type_name":"body","layer_1":{"zPos":10,"male":"body/human/male/"},"animations":["walk"]}`)
	addFile("spritesheets/.gitkeep", "")
	addFile("palette_definitions/"+marker+".txt", marker)

	if err := archive.Close(); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}
