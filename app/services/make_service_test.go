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
	required := []string{"schema_version", "command", "recipe", "assets", "output", "render", "layers", "warnings"}
	for _, key := range required {
		if _, ok := got[key]; !ok {
			t.Fatalf("missing report key %q in %+v", key, got)
		}
	}

	render := got["render"].(map[string]interface{})
	animationIDs := render["animation_ids"].([]interface{})
	if len(animationIDs) != 2 || animationIDs[0].(string) != "idle" || animationIDs[1].(string) != "walk" {
		t.Fatalf("unexpected animation ordering: %+v", animationIDs)
	}

	layers := got["layers"].(map[string]interface{})
	applied := layers["applied"].([]interface{})
	if len(applied) != 2 || applied[0].(string) != "body_human" || applied[1].(string) != "sword_training" {
		t.Fatalf("unexpected applied ordering: %+v", applied)
	}
}

func TestMakeServiceDeterministicPNGHashAndDimensions(t *testing.T) {
	assets, recipe := writeMakeFixture(t)
	out := filepath.Join(t.TempDir(), "sprite.png")

	result, problem := NewMakeService().Make(recipe, assets, out, "")
	if problem != nil {
		t.Fatalf("expected success, got %+v", problem)
	}
	if result.Summary.Canvas.Width != 8 || result.Summary.Canvas.Height != 8 {
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
	if img.Bounds().Dx() != 8 || img.Bounds().Dy() != 8 {
		t.Fatalf("unexpected dimensions: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
	}
	if got := fileSHA256(t, out); got != "08e2ea94d65b62e90bf48000bb9f1746ec09b3537d16f6d948b105d4ee4baa9e" {
		t.Fatalf("unexpected hash: %s", got)
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
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
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
