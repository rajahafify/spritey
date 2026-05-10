package services

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestMakeBatchServiceTwoJobSuccessInManifestOrder(t *testing.T) {
	assets, recipe := writeMakeFixture(t)
	root := t.TempDir()
	manifest := filepath.Join(root, "batch.json")
	jobAOut := filepath.Join(root, "out", "job-a.png")
	jobAReport := filepath.Join(root, "reports", "job-a.report.json")
	jobBOut := filepath.Join(root, "out", "job-b.png")

	writeFixtureFile(t, manifest, fmt.Sprintf(`{
  "schema_version": "1",
  "jobs": [
    {"id":"job-a","recipe":%q,"out":%q,"report":%q},
    {"id":"job-b","recipe":%q,"out":%q}
  ]
}`, recipe, jobAOut, jobAReport, recipe, jobBOut))

	result, problem := NewMakeBatchService().Make(manifest, assets)
	if problem != nil {
		t.Fatalf("expected success, got %+v", problem)
	}
	if result.Command != "make-batch" {
		t.Fatalf("unexpected command: %+v", result)
	}
	if result.Summary.JobCount != 2 || result.Summary.SuccessCount != 2 || result.Summary.FailedCount != 0 {
		t.Fatalf("unexpected summary: %+v", result.Summary)
	}
	if len(result.Jobs) != 2 || result.Jobs[0].ID != "job-a" || result.Jobs[1].ID != "job-b" {
		t.Fatalf("unexpected jobs order: %+v", result.Jobs)
	}
	if result.Jobs[0].Outputs.PNG.Path != jobAOut || result.Jobs[1].Outputs.PNG.Path != jobBOut {
		t.Fatalf("unexpected outputs: %+v", result.Jobs)
	}
}

func TestMakeBatchServiceResolvesRelativePathsAgainstManifestDirectory(t *testing.T) {
	assets, recipe := writeMakeFixture(t)
	root := t.TempDir()
	manifestDir := filepath.Join(root, "manifests")
	recipeDir := filepath.Join(manifestDir, "recipes")
	recipeCopy := filepath.Join(recipeDir, "copy.json")
	data, err := os.ReadFile(recipe)
	if err != nil {
		t.Fatal(err)
	}
	writeFixtureFile(t, recipeCopy, string(data))

	manifest := filepath.Join(manifestDir, "batch.json")
	writeFixtureFile(t, manifest, `{
  "schema_version":"1",
  "jobs":[
    {"id":"rel","recipe":"recipes/copy.json","out":"output/rel.png","report":"reports/rel.report.json"}
  ]
}`)

	result, problem := NewMakeBatchService().Make(manifest, assets)
	if problem != nil {
		t.Fatalf("expected success, got %+v", problem)
	}
	if len(result.Jobs) != 1 {
		t.Fatalf("expected one job, got %+v", result.Jobs)
	}
	if result.Jobs[0].Recipe != recipeCopy {
		t.Fatalf("expected resolved recipe path %q, got %q", recipeCopy, result.Jobs[0].Recipe)
	}
	if result.Jobs[0].Outputs.PNG.Path != filepath.Join(manifestDir, "output", "rel.png") {
		t.Fatalf("expected resolved out path, got %+v", result.Jobs[0].Outputs)
	}
	if result.Jobs[0].Outputs.Report == nil || result.Jobs[0].Outputs.Report.Path != filepath.Join(manifestDir, "reports", "rel.report.json") {
		t.Fatalf("expected resolved report path, got %+v", result.Jobs[0].Outputs)
	}
}

func TestMakeBatchServiceFailFastStopsAfterFirstFailure(t *testing.T) {
	assets, recipe := writeMakeFixture(t)
	root := t.TempDir()
	manifest := filepath.Join(root, "batch.json")
	wouldPassOut := filepath.Join(root, "out", "would-pass.png")
	writeFixtureFile(t, manifest, fmt.Sprintf(`{
  "schema_version":"1",
  "jobs":[
    {"id":"bad-first","recipe":%q,"out":%q},
    {"id":"would-pass","recipe":%q,"out":%q}
  ]
}`, filepath.Join(root, "missing.json"), filepath.Join(root, "out", "bad.png"), recipe, wouldPassOut))

	result, problem := NewMakeBatchService().Make(manifest, assets)
	if problem == nil {
		t.Fatal("expected batch failure")
	}
	if problem.Code != "RECIPE_FILE_NOT_FOUND" {
		t.Fatalf("unexpected failure code: %+v", problem)
	}
	if len(result.Jobs) != 1 || result.Jobs[0].ID != "bad-first" {
		t.Fatalf("expected fail-fast first job only, got %+v", result.Jobs)
	}
	if _, err := os.Stat(wouldPassOut); !os.IsNotExist(err) {
		t.Fatalf("expected second output not written due to fail-fast, stat err=%v", err)
	}
}
