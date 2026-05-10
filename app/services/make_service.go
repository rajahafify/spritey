package services

import (
	"encoding/json"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"sort"

	"github.com/rajahafify/spritey/app/models"
)

type MakeService struct {
	validator RecipeValidator
	loader    CatalogLoader
}

func NewMakeService() MakeService {
	return MakeService{
		validator: NewRecipeValidator(),
		loader:    NewCatalogLoader(),
	}
}

func (service MakeService) Make(recipePath string, assetsPath string, outPath string, reportPath string) (models.MakeResult, *models.MakeProblem) {
	validation, problem := service.validator.Validate(recipePath, assetsPath)
	if problem != nil {
		return models.MakeResult{}, makeProblemFromProblem(*problem)
	}

	catalog, problem := service.loader.Load(assetsPath)
	if problem != nil {
		return models.MakeResult{}, makeProblemFromProblem(*problem)
	}

	layerByID := indexLayers(catalog)
	appliedLayers := make([]indexedLayer, 0, len(validation.Selections))
	appliedIDs := make([]string, 0, len(validation.Selections))
	for _, selection := range validation.Selections {
		layer := layerByID[selection.ID]
		appliedLayers = append(appliedLayers, layer)
		appliedIDs = append(appliedIDs, selection.ID)
	}
	sort.Slice(appliedLayers, func(i, j int) bool {
		if appliedLayers[i].Layer.ZPos == appliedLayers[j].Layer.ZPos {
			return appliedLayers[i].Layer.ID < appliedLayers[j].Layer.ID
		}
		return appliedLayers[i].Layer.ZPos < appliedLayers[j].Layer.ZPos
	})

	animationIDs := make([]string, len(catalog.Pack.Defaults.Animations))
	copy(animationIDs, catalog.Pack.Defaults.Animations)
	sort.Strings(animationIDs)
	if len(animationIDs) == 0 {
		animationIDs = []string{"idle"}
	}

	canvasWidth := catalog.Pack.Defaults.CanvasWidth
	if canvasWidth <= 0 {
		canvasWidth = 64
	}
	canvasHeight := canvasWidth
	canvas := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))

	for _, layer := range appliedLayers {
		layerPathPrefix := layer.Layer.PathPrefix
		if layer.Layer.BodyTypePaths != nil {
			if matchedPath, ok := layer.Layer.BodyTypePaths[validation.BodyType]; ok && matchedPath != "" {
				layerPathPrefix = matchedPath
			}
		}
		sourcePath := filepath.Join(assetsPath, "spritesheets", filepath.FromSlash(layerPathPrefix), animationIDs[0]+".png")
		file, err := os.Open(sourcePath)
		if err != nil {
			return models.MakeResult{}, &models.MakeProblem{
				Code:    "RENDER_FAILED",
				Message: err.Error(),
				Field:   "render",
			}
		}
		src, err := png.Decode(file)
		file.Close()
		if err != nil {
			return models.MakeResult{}, &models.MakeProblem{
				Code:    "RENDER_FAILED",
				Message: err.Error(),
				Field:   "render",
			}
		}
		draw.Draw(canvas, canvas.Bounds(), src, src.Bounds().Min, draw.Over)
	}

	if err := writePNGAtomically(outPath, canvas); err != nil {
		return models.MakeResult{}, &models.MakeProblem{
			Code:    "RENDER_FAILED",
			Message: err.Error(),
			Field:   "out",
		}
	}

	result := models.MakeResult{
		Command: "make",
		Outputs: models.MakeOutputs{
			PNG: models.MakeFileOutput{Path: outPath},
		},
		Summary: models.MakeSummary{
			FrameCount: len(animationIDs),
			Canvas: models.MakeCanvas{
				Width:  canvasWidth,
				Height: canvasHeight,
			},
			AnimationCount: len(animationIDs),
		},
		Warnings: []string{},
	}

	if reportPath != "" {
		report := models.MakeReportV1{
			SchemaVersion: "1",
			Command:       "make",
			Warnings:      []string{},
		}
		report.Recipe.Path = recipePath
		report.Assets.Path = assetsPath
		report.Output.PNG.Path = outPath
		report.Render.Canvas.Width = canvasWidth
		report.Render.Canvas.Height = canvasHeight
		report.Render.FrameCount = len(animationIDs)
		report.Render.AnimationIDs = animationIDs
		report.Layers.Applied = appliedIDs

		if err := writeJSONFile(reportPath, report); err != nil {
			return models.MakeResult{}, &models.MakeProblem{
				Code:    "RENDER_FAILED",
				Message: err.Error(),
				Field:   "report",
			}
		}
		result.Outputs.Report = &models.MakeFileOutput{Path: reportPath}
	}

	return result, nil
}

func writePNGAtomically(path string, img image.Image) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tempFile, err := os.CreateTemp(filepath.Dir(path), "spritey-render-*.png")
	if err != nil {
		return err
	}
	tempName := tempFile.Name()
	removeTemp := true
	defer func() {
		if removeTemp {
			_ = os.Remove(tempName)
		}
	}()

	if err := png.Encode(tempFile, img); err != nil {
		tempFile.Close()
		return err
	}
	if err := tempFile.Close(); err != nil {
		return err
	}
	if err := os.Rename(tempName, path); err != nil {
		return err
	}
	removeTemp = false
	return nil
}

func writeJSONFile(path string, value interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func makeProblemFromProblem(problem models.Problem) *models.MakeProblem {
	return &models.MakeProblem{
		Code:    problem.Code,
		Message: problem.Message,
		Field:   problem.Field,
	}
}
