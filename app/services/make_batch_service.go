package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rajahafify/spritey/app/models"
)

type MakeBatchService struct {
	makeService MakeService
}

func NewMakeBatchService() MakeBatchService {
	return MakeBatchService{makeService: NewMakeService()}
}

func (service MakeBatchService) Make(manifestPath string, assetsPath string) (models.MakeBatchResult, *models.MakeProblem) {
	result := models.MakeBatchResult{
		Command:  "make-batch",
		Jobs:     []models.MakeBatchJobResult{},
		Warnings: []string{},
	}

	manifest, manifestDir, manifestProblem := loadBatchManifest(manifestPath)
	if manifestProblem != nil {
		return result, manifestProblem
	}

	result.Summary.JobCount = len(manifest.Jobs)
	for index, job := range manifest.Jobs {
		resolved := resolveBatchManifestJob(job, manifestDir)
		jobResult := models.MakeBatchJobResult{
			ID:       resolved.ID,
			Recipe:   resolved.Recipe,
			Outputs:  makeBatchOutputsForJob(resolved),
			Summary:  models.MakeSummary{},
			Warnings: []string{},
			Errors:   []models.MakeProblem{},
		}

		makeResult, makeProblem := service.makeService.Make(resolved.Recipe, assetsPath, resolved.Out, resolved.Report)
		if makeProblem != nil {
			jobResult.Errors = []models.MakeProblem{*makeProblem}
			result.Jobs = append(result.Jobs, jobResult)
			result.Summary.FailedCount = 1

			failed := *makeProblem
			details := map[string]interface{}{
				"job_id":                  resolved.ID,
				"job_index":               index + 1,
				"recipe":                  resolved.Recipe,
				"out":                     resolved.Out,
				"manifest":                manifestPath,
				"manifest_schema_version": manifest.SchemaVersion,
			}
			if resolved.Report != "" {
				details["report"] = resolved.Report
			}
			if failed.Details != nil {
				for key, value := range failed.Details {
					details[key] = value
				}
			}
			failed.Details = details
			return result, &failed
		}

		jobResult.Outputs = makeResult.Outputs
		jobResult.Summary = makeResult.Summary
		jobResult.Warnings = append([]string{}, makeResult.Warnings...)
		result.Warnings = append(result.Warnings, makeResult.Warnings...)
		result.Summary.SuccessCount++
		result.Jobs = append(result.Jobs, jobResult)
	}

	return result, nil
}

func loadBatchManifest(manifestPath string) (models.MakeBatchManifestV1, string, *models.MakeProblem) {
	data, readErr := os.ReadFile(manifestPath)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return models.MakeBatchManifestV1{}, "", &models.MakeProblem{
				Code:    "BATCH_MANIFEST_FILE_NOT_FOUND",
				Message: fmt.Sprintf("batch manifest file not found: %s", manifestPath),
				Field:   "manifest",
			}
		}
		return models.MakeBatchManifestV1{}, "", &models.MakeProblem{
			Code:    "READ_BATCH_MANIFEST_FAILED",
			Message: readErr.Error(),
			Field:   "manifest",
		}
	}

	var manifest models.MakeBatchManifestV1
	if unmarshalErr := json.Unmarshal(data, &manifest); unmarshalErr != nil {
		return models.MakeBatchManifestV1{}, "", &models.MakeProblem{
			Code:    "INVALID_BATCH_MANIFEST_JSON",
			Message: unmarshalErr.Error(),
			Field:   "manifest",
		}
	}

	if manifest.SchemaVersion != "1" {
		return models.MakeBatchManifestV1{}, "", &models.MakeProblem{
			Code:    "UNSUPPORTED_BATCH_MANIFEST_SCHEMA",
			Message: fmt.Sprintf("unsupported batch manifest schema_version: %s", manifest.SchemaVersion),
			Field:   "manifest.schema_version",
		}
	}

	if len(manifest.Jobs) == 0 {
		return models.MakeBatchManifestV1{}, "", &models.MakeProblem{
			Code:    "EMPTY_BATCH_JOBS",
			Message: "batch manifest jobs must not be empty",
			Field:   "manifest.jobs",
		}
	}

	for i, job := range manifest.Jobs {
		if job.ID == "" {
			return models.MakeBatchManifestV1{}, "", &models.MakeProblem{
				Code:    "MISSING_BATCH_JOB_ID",
				Message: fmt.Sprintf("batch job %d id is required", i+1),
				Field:   "manifest.jobs.id",
			}
		}
		if job.Recipe == "" {
			return models.MakeBatchManifestV1{}, "", &models.MakeProblem{
				Code:    "MISSING_BATCH_JOB_RECIPE",
				Message: fmt.Sprintf("batch job %d recipe is required", i+1),
				Field:   "manifest.jobs.recipe",
			}
		}
		if job.Out == "" {
			return models.MakeBatchManifestV1{}, "", &models.MakeProblem{
				Code:    "MISSING_BATCH_JOB_OUT",
				Message: fmt.Sprintf("batch job %d out is required", i+1),
				Field:   "manifest.jobs.out",
			}
		}
	}

	manifestDir := filepath.Dir(manifestPath)
	return manifest, manifestDir, nil
}

func resolveBatchManifestJob(job models.MakeBatchManifestJob, manifestDir string) models.MakeBatchManifestJob {
	resolved := job
	resolved.Recipe = resolveBatchManifestPath(manifestDir, job.Recipe)
	resolved.Out = resolveBatchManifestPath(manifestDir, job.Out)
	resolved.Report = resolveBatchManifestPath(manifestDir, job.Report)
	return resolved
}

func resolveBatchManifestPath(manifestDir string, value string) string {
	if value == "" {
		return ""
	}
	if filepath.IsAbs(value) {
		return filepath.Clean(value)
	}
	return filepath.Clean(filepath.Join(manifestDir, value))
}

func makeBatchOutputsForJob(job models.MakeBatchManifestJob) models.MakeOutputs {
	outputs := models.MakeOutputs{
		PNG: models.MakeFileOutput{Path: job.Out},
	}
	if job.Report != "" {
		outputs.Report = &models.MakeFileOutput{Path: job.Report}
	}
	return outputs
}
