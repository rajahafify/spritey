package services

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rajahafify/spritey/app/models"
)

type MakeService struct {
	validator     RecipeValidator
	loader        CatalogLoader
	frameResolver SpriteFrameResolver
	recolorer     PaletteRecolorer
}

var outputPNGArtifactFn = computeOutputPNGArtifact

const (
	lpcStripWidth      = 832
	lpcFallbackHeight  = 256
)

var lpcAnimationOrder = []string{"spellcast", "thrust", "walk", "slash", "shoot", "hurt"}

func NewMakeService() MakeService {
	return MakeService{
		validator:     NewRecipeValidator(),
		loader:        NewCatalogLoader(),
		frameResolver: NewSpriteFrameResolver(),
		recolorer:     NewPaletteRecolorer(),
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
	renderInputByLayerID := map[string]models.RecipeRenderInput{}
	for _, input := range validation.RenderInputs {
		renderInputByLayerID[input.LayerID] = input
	}

	type appliedLayer struct {
		Indexed indexedLayer
		Render  models.RecipeRenderInput
		Palette string
	}
	appliedLayers := make([]appliedLayer, 0, len(validation.Selections))
	appliedIDs := make([]string, 0, len(validation.Selections))
	for _, selection := range validation.Selections {
		layer := layerByID[selection.ID]
		appliedLayers = append(appliedLayers, appliedLayer{
			Indexed: layer,
			Render:  renderInputByLayerID[selection.ID],
			Palette: selection.PaletteVariant,
		})
		appliedIDs = append(appliedIDs, selection.ID)
	}
	sort.Slice(appliedLayers, func(i, j int) bool {
		if appliedLayers[i].Indexed.Layer.ZPos == appliedLayers[j].Indexed.Layer.ZPos {
			return appliedLayers[i].Indexed.Layer.ID < appliedLayers[j].Indexed.Layer.ID
		}
		return appliedLayers[i].Indexed.Layer.ZPos < appliedLayers[j].Indexed.Layer.ZPos
	})

	type emittedRow struct {
		AnimationID string
		Image       *image.RGBA
	}

	emittedRows := []emittedRow{}
	for _, animationID := range lpcAnimationOrder {
		var row *image.RGBA
		rowHeight := 0

		for _, layer := range appliedLayers {
			resolvedFramePath, found, err := service.frameResolver.ResolveFrame(assetsPath, layer.Render.ResolvedPath, animationID)
			if err != nil {
				return models.MakeResult{}, &models.MakeProblem{
					Code:    "RENDER_FAILED",
					Message: err.Error(),
					Field:   "render",
				}
			}
			if !found {
				continue
			}

			sourcePath := filepath.Join(assetsPath, "spritesheets", filepath.FromSlash(resolvedFramePath))
			src, err := loadPNG(sourcePath)
			if err != nil {
				return models.MakeResult{}, &models.MakeProblem{
					Code:    "RENDER_FAILED",
					Message: err.Error(),
					Field:   "render",
				}
			}
			if layer.Indexed.Layer.RecolorMaterial != "" && layer.Palette != "" {
				src = service.recolorer.Recolor(src, assetsPath, layer.Indexed.Layer.RecolorMaterial, layer.Palette)
			}

			if row == nil {
				rowHeight = src.Bounds().Dy()
				row = padLayerToRow(src, rowHeight)
				continue
			}

			draw.Draw(row, row.Bounds(), padLayerToRow(src, rowHeight), image.Point{}, draw.Over)
		}

		if row != nil {
			emittedRows = append(emittedRows, emittedRow{
				AnimationID: animationID,
				Image:       row,
			})
		}
	}

	canvasWidth := lpcStripWidth
	canvasHeight := lpcFallbackHeight
	animationIDs := []string{}
	if len(emittedRows) > 0 {
		canvasHeight = 0
		for _, row := range emittedRows {
			canvasHeight += row.Image.Bounds().Dy()
			animationIDs = append(animationIDs, row.AnimationID)
		}
	}
	canvas := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))
	if len(emittedRows) > 0 {
		offsetY := 0
		for _, row := range emittedRows {
			draw.Draw(canvas, row.Image.Bounds().Add(image.Pt(0, offsetY)), row.Image, row.Image.Bounds().Min, draw.Over)
			offsetY += row.Image.Bounds().Dy()
		}
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
		Warnings: append([]string{}, validation.Warnings...),
	}

	if reportPath != "" {
		outputArtifact, err := outputPNGArtifactFn(outPath)
		if err != nil {
			return models.MakeResult{}, &models.MakeProblem{
				Code:    "RENDER_FAILED",
				Message: err.Error(),
				Field:   "report",
			}
		}

		composedLayers := make([]models.MakeReportComposedLayer, 0, len(appliedLayers))
		for _, layer := range appliedLayers {
			composedLayers = append(composedLayers, models.MakeReportComposedLayer{
				Category:         layer.Indexed.Category,
				ID:               layer.Indexed.Layer.ID,
				ZPos:             layer.Indexed.Layer.ZPos,
				ResolvedBodyType: layer.Render.ResolvedBodyType,
				ResolvedPath:     normalizeSlashPath(layer.Render.ResolvedPath),
				PaletteVariant:   layer.Palette,
				Credits:          copyCredits(layer.Indexed.Layer.Credits),
			})
		}

		report := models.MakeReportV1Provenance{
			SchemaVersion: "1",
			Command:       "make",
			Warnings:      append([]string{}, validation.Warnings...),
		}
		report.Pack.ID = catalog.Pack.ID
		report.Pack.Name = catalog.Pack.Name
		report.Recipe.Path = recipePath
		report.Recipe.BodyTypeEffective = validation.BodyType
		report.Recipe.BodyTypeRequested = validation.BodyTypeRequested
		report.Assets.Path = assetsPath
		report.Output.PNG.Path = outPath
		report.Artifacts.OutputPNG = outputArtifact
		report.Render.Canvas.Width = canvasWidth
		report.Render.Canvas.Height = canvasHeight
		report.Render.FrameCount = len(animationIDs)
		report.Render.AnimationIDs = animationIDs
		report.Layers.Applied = appliedIDs
		report.Layers.Composed = composedLayers

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

func loadPNG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return png.Decode(file)
}

func padLayerToRow(src image.Image, rowHeight int) *image.RGBA {
	if rowHeight <= 0 {
		rowHeight = src.Bounds().Dy()
	}
	row := image.NewRGBA(image.Rect(0, 0, lpcStripWidth, rowHeight))
	draw.Draw(row, src.Bounds(), src, src.Bounds().Min, draw.Src)
	return row
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

func computeOutputPNGArtifact(path string) (models.MakeReportOutputPNGArtifact, error) {
	file, err := os.Open(path)
	if err != nil {
		return models.MakeReportOutputPNGArtifact{}, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return models.MakeReportOutputPNGArtifact{}, err
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return models.MakeReportOutputPNGArtifact{}, err
	}

	return models.MakeReportOutputPNGArtifact{
		SHA256: fmt.Sprintf("%x", hash.Sum(nil)),
		Bytes:  info.Size(),
	}, nil
}

func makeProblemFromProblem(problem models.Problem) *models.MakeProblem {
	return &models.MakeProblem{
		Code:    problem.Code,
		Message: problem.Message,
		Field:   problem.Field,
	}
}

func normalizeSlashPath(raw string) string {
	normalized := strings.ReplaceAll(raw, "\\", "/")
	return path.Clean(normalized)
}

func copyCredits(credits []models.Credit) []models.Credit {
	if len(credits) == 0 {
		return []models.Credit{}
	}
	cloned := make([]models.Credit, len(credits))
	copy(cloned, credits)
	return cloned
}
