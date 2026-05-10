package services

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/rajahafify/spritey/app/models"
)

func TestMakeServiceReportV1FieldsAndOrdering(t *testing.T) {
	assets, recipe := writeMakeFixture(t)
	out := filepath.Join(t.TempDir(), "sprite.png")
	report := filepath.Join(t.TempDir(), "sprite.report.json")

	result, problem := NewMakeService().Make(recipe, assets, out, report)
	if problem != nil {
		t.Fatalf("expected success, got %+v", problem)
	}
	if result.Outputs.Report == nil {
		t.Fatal("expected report output")
	}

	reportData, err := os.ReadFile(report)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(reportData, &got); err != nil {
		t.Fatal(err)
	}
	required := []string{"schema_version", "command", "pack", "recipe", "assets", "output", "artifacts", "render", "layers", "warnings"}
	for _, key := range required {
		if _, ok := got[key]; !ok {
			t.Fatalf("missing report key %q in %+v", key, got)
		}
	}
	pack := got["pack"].(map[string]interface{})
	if pack["id"] != "make-assets" || pack["name"] != "Make Assets" {
		t.Fatalf("unexpected pack metadata: %+v", pack)
	}
	recipeMeta := got["recipe"].(map[string]interface{})
	if recipeMeta["path"] != recipe {
		t.Fatalf("unexpected recipe path: %+v", recipeMeta)
	}
	if recipeMeta["body_type_requested"] != "male" || recipeMeta["body_type_effective"] != "male" {
		t.Fatalf("unexpected recipe body type provenance: %+v", recipeMeta)
	}
	artifacts := got["artifacts"].(map[string]interface{})
	outputPNG := artifacts["output_png"].(map[string]interface{})
	if outputPNG["sha256"] != fileSHA256(t, out) {
		t.Fatalf("unexpected report output sha256: %+v", outputPNG)
	}
	info, err := os.Stat(out)
	if err != nil {
		t.Fatal(err)
	}
	if int64(outputPNG["bytes"].(float64)) != info.Size() {
		t.Fatalf("unexpected report output bytes: got=%v want=%d", outputPNG["bytes"], info.Size())
	}

	render := got["render"].(map[string]interface{})
	animationIDs := render["animation_ids"].([]interface{})
	if len(animationIDs) != 1 || animationIDs[0].(string) != "walk" {
		t.Fatalf("unexpected animation ordering: %+v", animationIDs)
	}
	if int(render["frame_count"].(float64)) != 1 {
		t.Fatalf("unexpected report frame_count: %+v", render["frame_count"])
	}

	layers := got["layers"].(map[string]interface{})
	applied := layers["applied"].([]interface{})
	if len(applied) != 2 || applied[0].(string) != "body_human" || applied[1].(string) != "sword_training" {
		t.Fatalf("unexpected applied ordering: %+v", applied)
	}
	composed := layers["composed"].([]interface{})
	if len(composed) != 2 {
		t.Fatalf("unexpected composed length: %+v", composed)
	}
	assertComposedLayer(t, composed[0], "body", "body_human", 10, "male", "body/human/male", "", 0)
	assertComposedLayer(t, composed[1], "weapon", "sword_training", 30, "male", "weapon/sword/male", "", 0)
}

func TestMakeServiceDeterministicPNGHashAndDimensions(t *testing.T) {
	assets, recipe := writeMakeFixture(t)
	out := filepath.Join(t.TempDir(), "sprite.png")

	result, problem := NewMakeService().Make(recipe, assets, out, "")
	if problem != nil {
		t.Fatalf("expected success, got %+v", problem)
	}
	if result.Summary.Canvas.Width != 832 || result.Summary.Canvas.Height != 8 {
		t.Fatalf("unexpected canvas size: %+v", result.Summary.Canvas)
	}

	file, err := os.Open(out)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != 832 || img.Bounds().Dy() != 8 {
		t.Fatalf("unexpected dimensions: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
	}
	if got := fileSHA256(t, out); got != "2df8e9b2d6cd1b5c732cead7598e355558517ad667b6cec34721c4019c0c268a" {
		t.Fatalf("unexpected hash: %s", got)
	}
}

func TestMakeServiceUsesFixedLPCAnimationOrderAndSkipsMissingLayerFrames(t *testing.T) {
	assets, recipe := writeMakeLPCParityFixture(t, makeLPCParityFixtureOptions{
		requiredAnimations: []string{"idle", "walk"},
		bodyWalkSize:       image.Pt(8, 8),
		weaponWalkSize:     image.Pt(8, 8),
		includeBodySlash:   true,
		includeWeaponSlash: false,
	})
	out := filepath.Join(t.TempDir(), "sprite.png")
	report := filepath.Join(t.TempDir(), "sprite.report.json")

	result, problem := NewMakeService().Make(recipe, assets, out, report)
	if problem != nil {
		t.Fatalf("expected success, got %+v", problem)
	}
	if result.Summary.FrameCount != 2 || result.Summary.AnimationCount != 2 {
		t.Fatalf("expected two emitted rows, got %+v", result.Summary)
	}

	file, err := os.Open(out)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != 832 || img.Bounds().Dy() != 16 {
		t.Fatalf("unexpected dimensions: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
	}

	top := color.RGBAModel.Convert(img.At(2, 2)).(color.RGBA)
	next := color.RGBAModel.Convert(img.At(2, 10)).(color.RGBA)
	if top != (color.RGBA{R: 214, G: 50, B: 90, A: 255}) {
		t.Fatalf("expected first emitted row to be walk overlay, got %+v", top)
	}
	if next != (color.RGBA{R: 74, G: 130, B: 180, A: 255}) {
		t.Fatalf("expected second emitted row to be body slash, got %+v", next)
	}

	reportData, err := os.ReadFile(report)
	if err != nil {
		t.Fatal(err)
	}
	var reportJSON map[string]interface{}
	if err := json.Unmarshal(reportData, &reportJSON); err != nil {
		t.Fatal(err)
	}
	render := reportJSON["render"].(map[string]interface{})
	animationIDs := render["animation_ids"].([]interface{})
	if len(animationIDs) != 2 || animationIDs[0] != "walk" || animationIDs[1] != "slash" {
		t.Fatalf("unexpected emitted animation ids: %+v", animationIDs)
	}
	if int(render["frame_count"].(float64)) != 2 {
		t.Fatalf("unexpected emitted frame_count: %+v", render["frame_count"])
	}
}

func TestMakeServiceRowHeightUsesFirstContributingLayerAndPadsSubsequentLayers(t *testing.T) {
	assets, recipe := writeMakeLPCParityFixture(t, makeLPCParityFixtureOptions{
		requiredAnimations: []string{"idle", "walk"},
		bodyWalkSize:       image.Pt(8, 6),
		weaponWalkSize:     image.Pt(8, 10),
	})
	out := filepath.Join(t.TempDir(), "sprite.png")

	result, problem := NewMakeService().Make(recipe, assets, out, "")
	if problem != nil {
		t.Fatalf("expected success, got %+v", problem)
	}
	if result.Summary.FrameCount != 1 {
		t.Fatalf("expected single emitted row, got %d", result.Summary.FrameCount)
	}
	if result.Summary.Canvas.Width != 832 || result.Summary.Canvas.Height != 6 {
		t.Fatalf("unexpected canvas size for first-layer row height parity: %+v", result.Summary.Canvas)
	}

	file, err := os.Open(out)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != 832 || img.Bounds().Dy() != 6 {
		t.Fatalf("unexpected rendered dimensions: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestMakeServiceNoEmittedRowsReturnsTransparentFallbackCanvas(t *testing.T) {
	assets, recipe := writeMakeReadinessFixture(t, makeReadinessFixtureOptions{
		recipeBodyType:      "male",
		missingFallback:     "male",
		requiredAnimations:  []string{"idle"},
		includeBodyIdle:     true,
		includeWeaponIdle:   true,
		weaponHasFemalePath: true,
	})
	out := filepath.Join(t.TempDir(), "sprite.png")

	result, problem := NewMakeService().Make(recipe, assets, out, "")
	if problem != nil {
		t.Fatalf("expected success, got %+v", problem)
	}
	if result.Summary.FrameCount != 0 || result.Summary.AnimationCount != 0 {
		t.Fatalf("expected no emitted rows in summary, got %+v", result.Summary)
	}
	if result.Summary.Canvas.Width != 832 || result.Summary.Canvas.Height != 256 {
		t.Fatalf("unexpected fallback canvas: %+v", result.Summary.Canvas)
	}

	file, err := os.Open(out)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	if img.Bounds().Dx() != 832 || img.Bounds().Dy() != 256 {
		t.Fatalf("unexpected fallback png dimensions: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestMakeServiceMissingRequiredFrameReturnsInputValidationAndNoOutput(t *testing.T) {
	assets, recipe := writeMakeReadinessFixture(t, makeReadinessFixtureOptions{
		recipeBodyType:      "male",
		missingFallback:     "male",
		requiredAnimations:  []string{"idle", "walk"},
		includeBodyIdle:     true,
		includeBodyWalk:     true,
		includeWeaponIdle:   true,
		includeWeaponWalk:   false,
		weaponHasFemalePath: true,
	})
	out := filepath.Join(t.TempDir(), "sprite.png")

	_, problem := NewMakeService().Make(recipe, assets, out, "")
	if problem == nil {
		t.Fatal("expected make problem")
	}
	if problem.Code != "MISSING_SPRITE_FRAME" {
		t.Fatalf("unexpected make code: %+v", problem)
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Fatalf("expected no output file, stat err=%v", err)
	}
}

func TestMakeServiceResolvesMappedAttackSlashPath(t *testing.T) {
	assets, recipe := writeSlashParityFixture(t, true)
	out := filepath.Join(t.TempDir(), "sprite.png")

	result, problem := NewMakeService().Make(recipe, assets, out, "")
	if problem != nil {
		t.Fatalf("expected success, got %+v", problem)
	}
	if result.Summary.FrameCount != 1 {
		t.Fatalf("expected one frame, got %+v", result.Summary)
	}
}

func TestMakeServiceMappedAttackSlashMissingReturnsMissingSpriteFrame(t *testing.T) {
	assets, recipe := writeSlashParityFixture(t, false)
	out := filepath.Join(t.TempDir(), "sprite.png")

	_, problem := NewMakeService().Make(recipe, assets, out, "")
	if problem == nil {
		t.Fatal("expected make problem")
	}
	if problem.Code != "MISSING_SPRITE_FRAME" {
		t.Fatalf("unexpected make code: %+v", problem)
	}
}

func TestMakeServiceFallbackPathWarningsInResultAndReport(t *testing.T) {
	assets, recipe := writeMakeReadinessFixture(t, makeReadinessFixtureOptions{
		recipeBodyType:        "female",
		missingFallback:       "male",
		requiredAnimations:    []string{"idle"},
		includeBodyIdle:       true,
		includeBodyFemaleIdle: true,
		includeWeaponIdle:     true,
		weaponHasFemalePath:   false,
	})
	out := filepath.Join(t.TempDir(), "sprite.png")
	report := filepath.Join(t.TempDir(), "sprite.report.json")

	result, problem := NewMakeService().Make(recipe, assets, out, report)
	if problem != nil {
		t.Fatalf("expected make success, got %+v", problem)
	}
	if len(result.Warnings) != 1 {
		t.Fatalf("expected one warning in result, got %+v", result.Warnings)
	}

	reportData, err := os.ReadFile(report)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(reportData, &got); err != nil {
		t.Fatal(err)
	}
	warnings, ok := got["warnings"].([]interface{})
	if !ok || len(warnings) != 1 {
		t.Fatalf("expected one warning in report, got %+v", got["warnings"])
	}
	recipeMeta := got["recipe"].(map[string]interface{})
	if recipeMeta["body_type_requested"] != "female" || recipeMeta["body_type_effective"] != "female" {
		t.Fatalf("unexpected recipe body provenance: %+v", recipeMeta)
	}
	composed := got["layers"].(map[string]interface{})["composed"].([]interface{})
	if len(composed) != 2 {
		t.Fatalf("unexpected composed length: %+v", composed)
	}
	assertComposedLayer(t, composed[0], "body", "body_human", 10, "female", "body/human/female", "", 0)
	assertComposedLayer(t, composed[1], "weapon", "sword_training", 30, "male", "weapon/sword/male", "", 0)
}

func TestMakeServiceReportArtifactMetadataFailureReturnsRenderFailedAndSkipsReport(t *testing.T) {
	assets, recipe := writeMakeFixture(t)
	out := filepath.Join(t.TempDir(), "sprite.png")
	report := filepath.Join(t.TempDir(), "sprite.report.json")

	original := outputPNGArtifactFn
	outputPNGArtifactFn = func(path string) (models.MakeReportOutputPNGArtifact, error) {
		return models.MakeReportOutputPNGArtifact{}, fmt.Errorf("metadata failure")
	}
	defer func() {
		outputPNGArtifactFn = original
	}()

	_, problem := NewMakeService().Make(recipe, assets, out, report)
	if problem == nil {
		t.Fatal("expected make problem")
	}
	if problem.Code != "RENDER_FAILED" {
		t.Fatalf("unexpected make problem code: %+v", problem)
	}
	if problem.Field != "report" {
		t.Fatalf("unexpected make problem field: %+v", problem)
	}
	if _, err := os.Stat(report); !os.IsNotExist(err) {
		t.Fatalf("expected report file not written, stat err=%v", err)
	}
}

func assertComposedLayer(t *testing.T, value interface{}, category string, id string, zPos int, resolvedBodyType string, resolvedPath string, paletteVariant string, creditCount int) {
	t.Helper()
	entry := value.(map[string]interface{})
	if entry["category"] != category || entry["id"] != id {
		t.Fatalf("unexpected composed identity: %+v", entry)
	}
	if int(entry["z_pos"].(float64)) != zPos {
		t.Fatalf("unexpected composed z_pos: %+v", entry)
	}
	if entry["resolved_body_type"] != resolvedBodyType || entry["resolved_path"] != resolvedPath {
		t.Fatalf("unexpected resolved fields: %+v", entry)
	}
	if entry["palette_variant"] != paletteVariant {
		t.Fatalf("unexpected palette_variant: %+v", entry)
	}
	credits, ok := entry["credits"].([]interface{})
	if !ok {
		t.Fatalf("expected non-null credits array, got %+v", entry["credits"])
	}
	if len(credits) != creditCount {
		t.Fatalf("unexpected credits length: %+v", entry)
	}
}

func writeMakeFixture(t *testing.T) (string, string) {
	t.Helper()

	root := t.TempDir()
	assets := filepath.Join(root, "assets")
	writeFixtureFile(t, filepath.Join(assets, "pack.json"), `{"schema_version":"1","id":"make-assets","name":"Make Assets","defaults":{"body_type":"male","animations":["walk","idle"],"canvas_width":8}}`)
	writeFixtureFile(t, filepath.Join(assets, "sheet_definitions", "body", "body_human.json"), `{"name":"Human Body","type_name":"body","layer_1":{"zPos":10,"male":"body/human/male/"},"animations":["walk","idle"]}`)
	writeFixtureFile(t, filepath.Join(assets, "sheet_definitions", "weapon", "sword_training.json"), `{"name":"Training Sword","type_name":"weapon","layer_1":{"zPos":30,"male":"weapon/sword/male/"},"animations":["walk","idle"]}`)
	writeFixtureFile(t, filepath.Join(assets, "palette_definitions", ".gitkeep"), "")
	writeFixtureFile(t, filepath.Join(assets, "spritesheets", ".gitkeep"), "")

	writeLayerPNG(t, filepath.Join(assets, "spritesheets", "body", "human", "male", "idle.png"), color.RGBA{R: 40, G: 100, B: 200, A: 255})
	writeLayerPNG(t, filepath.Join(assets, "spritesheets", "body", "human", "male", "walk.png"), color.RGBA{R: 50, G: 110, B: 210, A: 255})
	writeLayerPNG(t, filepath.Join(assets, "spritesheets", "weapon", "sword", "male", "idle.png"), color.RGBA{R: 200, G: 40, B: 80, A: 255})
	writeLayerPNG(t, filepath.Join(assets, "spritesheets", "weapon", "sword", "male", "walk.png"), color.RGBA{R: 210, G: 50, B: 90, A: 255})

	recipe := filepath.Join(root, "recipe.json")
	writeFixtureFile(t, recipe, `{"body_type":"male","selections":{"body":{"id":"body_human"},"weapon":{"id":"sword_training"}}}`)
	return assets, recipe
}

type makeLPCParityFixtureOptions struct {
	requiredAnimations []string
	bodyWalkSize       image.Point
	weaponWalkSize     image.Point
	includeBodySlash   bool
	includeWeaponSlash bool
}

func writeMakeLPCParityFixture(t *testing.T, opts makeLPCParityFixtureOptions) (string, string) {
	t.Helper()

	requiredAnimations := opts.requiredAnimations
	if len(requiredAnimations) == 0 {
		requiredAnimations = []string{"idle", "walk"}
	}
	bodyWalkSize := opts.bodyWalkSize
	if bodyWalkSize.X <= 0 || bodyWalkSize.Y <= 0 {
		bodyWalkSize = image.Pt(8, 8)
	}
	weaponWalkSize := opts.weaponWalkSize
	if weaponWalkSize.X <= 0 || weaponWalkSize.Y <= 0 {
		weaponWalkSize = image.Pt(8, 8)
	}

	animationsJSON := `["` + requiredAnimations[0] + `"`
	for i := 1; i < len(requiredAnimations); i++ {
		animationsJSON += `,"` + requiredAnimations[i] + `"`
	}
	animationsJSON += `]`

	root := t.TempDir()
	assets := filepath.Join(root, "assets")
	writeFixtureFile(t, filepath.Join(assets, "pack.json"), `{"schema_version":"1","id":"make-lpc-parity","name":"Make LPC Parity","defaults":{"body_type":"male","animations":`+animationsJSON+`,"canvas_width":8}}`)
	writeFixtureFile(t, filepath.Join(assets, "sheet_definitions", "body", "body_human.json"), `{"name":"Human Body","type_name":"body","layer_1":{"zPos":10,"male":"body/human/male/"},"animations":["idle","walk","slash"]}`)
	writeFixtureFile(t, filepath.Join(assets, "sheet_definitions", "weapon", "sword_training.json"), `{"name":"Training Sword","type_name":"weapon","layer_1":{"zPos":30,"male":"weapon/sword/male/"},"animations":["idle","walk","slash"]}`)
	writeFixtureFile(t, filepath.Join(assets, "palette_definitions", ".gitkeep"), "")
	writeFixtureFile(t, filepath.Join(assets, "spritesheets", ".gitkeep"), "")

	writeLayerPNG(t, filepath.Join(assets, "spritesheets", "body", "human", "male", "idle.png"), color.RGBA{R: 40, G: 100, B: 200, A: 255})
	writeLayerPNG(t, filepath.Join(assets, "spritesheets", "weapon", "sword", "male", "idle.png"), color.RGBA{R: 200, G: 40, B: 80, A: 255})
	writeLayerPNGWithSize(t, filepath.Join(assets, "spritesheets", "body", "human", "male", "walk.png"), color.RGBA{R: 50, G: 110, B: 210, A: 255}, bodyWalkSize.X, bodyWalkSize.Y)
	writeLayerPNGWithSize(t, filepath.Join(assets, "spritesheets", "weapon", "sword", "male", "walk.png"), color.RGBA{R: 210, G: 50, B: 90, A: 255}, weaponWalkSize.X, weaponWalkSize.Y)
	if opts.includeBodySlash {
		writeLayerPNG(t, filepath.Join(assets, "spritesheets", "body", "human", "male", "slash.png"), color.RGBA{R: 70, G: 130, B: 180, A: 255})
	}
	if opts.includeWeaponSlash {
		writeLayerPNG(t, filepath.Join(assets, "spritesheets", "weapon", "sword", "male", "slash.png"), color.RGBA{R: 180, G: 60, B: 100, A: 255})
	}

	recipe := filepath.Join(root, "recipe.json")
	writeFixtureFile(t, recipe, `{"body_type":"male","selections":{"body":{"id":"body_human"},"weapon":{"id":"sword_training"}}}`)
	return assets, recipe
}

type makeReadinessFixtureOptions struct {
	recipeBodyType        string
	missingFallback       string
	requiredAnimations    []string
	includeBodyIdle       bool
	includeBodyFemaleIdle bool
	includeBodyWalk       bool
	includeWeaponIdle     bool
	includeWeaponWalk     bool
	weaponHasFemalePath   bool
}

func writeMakeReadinessFixture(t *testing.T, opts makeReadinessFixtureOptions) (string, string) {
	t.Helper()
	root := t.TempDir()
	assets := filepath.Join(root, "assets")

	recipeBodyType := opts.recipeBodyType
	if recipeBodyType == "" {
		recipeBodyType = "male"
	}
	missingFallback := opts.missingFallback
	if missingFallback == "" {
		missingFallback = "male"
	}
	requiredAnimations := opts.requiredAnimations
	if len(requiredAnimations) == 0 {
		requiredAnimations = []string{"idle"}
	}
	animationsJSON := `["` + requiredAnimations[0] + `"`
	for i := 1; i < len(requiredAnimations); i++ {
		animationsJSON += `,"` + requiredAnimations[i] + `"`
	}
	animationsJSON += `]`

	weaponFemale := `"female":"weapon/sword/female/"`
	if !opts.weaponHasFemalePath {
		weaponFemale = ""
	}

	writeFixtureFile(t, filepath.Join(assets, "pack.json"), `{"schema_version":"1","id":"make-readiness","name":"Make Readiness","defaults":{"body_type":"male","animations":`+animationsJSON+`,"canvas_width":8,"missing_body_type_fallback":"`+missingFallback+`"}}`)
	writeFixtureFile(t, filepath.Join(assets, "sheet_definitions", "body", "body_human.json"), `{"name":"Human Body","type_name":"body","layer_1":{"zPos":10,"male":"body/human/male/","female":"body/human/female/"},"animations":["idle","walk"]}`)
	writeFixtureFile(t, filepath.Join(assets, "sheet_definitions", "weapon", "sword_training.json"), `{"name":"Training Sword","type_name":"weapon","layer_1":{"zPos":30,"male":"weapon/sword/male/"`+withOptionalLeadingComma(weaponFemale)+`},"animations":["idle","walk"]}`)
	writeFixtureFile(t, filepath.Join(assets, "palette_definitions", ".gitkeep"), "")
	writeFixtureFile(t, filepath.Join(assets, "spritesheets", ".gitkeep"), "")

	if opts.includeBodyIdle {
		writeLayerPNG(t, filepath.Join(assets, "spritesheets", "body", "human", "male", "idle.png"), color.RGBA{R: 40, G: 100, B: 200, A: 255})
	}
	if opts.includeBodyWalk {
		writeLayerPNG(t, filepath.Join(assets, "spritesheets", "body", "human", "male", "walk.png"), color.RGBA{R: 50, G: 110, B: 210, A: 255})
	}
	if opts.includeBodyFemaleIdle {
		writeLayerPNG(t, filepath.Join(assets, "spritesheets", "body", "human", "female", "idle.png"), color.RGBA{R: 45, G: 105, B: 205, A: 255})
	}
	if opts.includeWeaponIdle {
		writeLayerPNG(t, filepath.Join(assets, "spritesheets", "weapon", "sword", "male", "idle.png"), color.RGBA{R: 200, G: 40, B: 80, A: 255})
	}
	if opts.includeWeaponWalk {
		writeLayerPNG(t, filepath.Join(assets, "spritesheets", "weapon", "sword", "male", "walk.png"), color.RGBA{R: 210, G: 50, B: 90, A: 255})
	}

	recipe := filepath.Join(root, "recipe.json")
	writeFixtureFile(t, recipe, `{"body_type":"`+recipeBodyType+`","selections":{"body":{"id":"body_human"},"weapon":{"id":"sword_training"}}}`)
	return assets, recipe
}

func withOptionalLeadingComma(value string) string {
	if value == "" {
		return ""
	}
	return "," + value
}

func writeSlashParityFixture(t *testing.T, includeMappedSlash bool) (string, string) {
	t.Helper()
	root := t.TempDir()
	assets := filepath.Join(root, "assets")

	writeFixtureFile(t, filepath.Join(assets, "pack.json"), `{"schema_version":"1","id":"slash-parity","name":"Slash Parity","defaults":{"body_type":"male","animations":["slash"],"canvas_width":8}}`)
	writeFixtureFile(t, filepath.Join(assets, "sheet_definitions", "body", "body_human.json"), `{"name":"Human Body","type_name":"body","layer_1":{"zPos":10,"male":"body/human/male/"},"animations":["slash"]}`)
	writeFixtureFile(t, filepath.Join(assets, "sheet_definitions", "weapon", "sword_training.json"), `{"name":"Training Sword","type_name":"weapon","layer_1":{"zPos":30,"male":"weapon/sword/male/"},"animations":["slash"]}`)
	writeFixtureFile(t, filepath.Join(assets, "palette_definitions", ".gitkeep"), "")
	writeFixtureFile(t, filepath.Join(assets, "spritesheets", ".gitkeep"), "")

	writeLayerPNG(t, filepath.Join(assets, "spritesheets", "body", "human", "male", "slash.png"), color.RGBA{R: 40, G: 100, B: 200, A: 255})
	if includeMappedSlash {
		writeLayerPNG(t, filepath.Join(assets, "spritesheets", "weapon", "sword", "male", "attack_slash", "front.png"), color.RGBA{R: 200, G: 40, B: 80, A: 255})
		writeLayerPNG(t, filepath.Join(assets, "spritesheets", "weapon", "sword", "male", "attack_slash", "behind_pose.png"), color.RGBA{R: 180, G: 20, B: 60, A: 255})
	}

	recipe := filepath.Join(root, "recipe.json")
	writeFixtureFile(t, recipe, `{"body_type":"male","selections":{"body":{"id":"body_human"},"weapon":{"id":"sword_training"}}}`)
	return assets, recipe
}

func writeFixtureFile(t *testing.T, path, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeLayerPNG(t *testing.T, path string, fill color.RGBA) {
	t.Helper()
	writeLayerPNGWithSize(t, path, fill, 8, 8)
}

func writeLayerPNGWithSize(t *testing.T, path string, fill color.RGBA, width int, height int) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := fill
			c.R = uint8((int(c.R) + x + y) % 255)
			img.SetRGBA(x, y, c)
		}
	}
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		t.Fatal(err)
	}
}

func fileSHA256(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	return fmt.Sprintf("%x", sum[:])
}
