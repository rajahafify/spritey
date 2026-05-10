package services

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

func TestPaletteRecolorerAppliesExactVariant(t *testing.T) {
	assets := t.TempDir()
	writePaletteFile(t, filepath.Join(assets, "palette_definitions", "skin_main.json"), `{"light":["#112233"],"tan":["#445566"]}`)

	input := image.NewRGBA(image.Rect(0, 0, 1, 1))
	input.SetRGBA(0, 0, color.RGBA{R: 0x11, G: 0x22, B: 0x33, A: 0xFF})

	out := NewPaletteRecolorer().Recolor(input, assets, "skin", "tan")
	got := color.RGBAModel.Convert(out.At(0, 0)).(color.RGBA)
	want := color.RGBA{R: 0x44, G: 0x55, B: 0x66, A: 0xFF}
	if got != want {
		t.Fatalf("unexpected recolor result: got=%+v want=%+v", got, want)
	}
}

func TestPaletteRecolorerUnknownVariantReturnsOriginal(t *testing.T) {
	assets := t.TempDir()
	writePaletteFile(t, filepath.Join(assets, "palette_definitions", "skin_main.json"), `{"light":["#112233"],"tan":["#445566"]}`)

	input := image.NewRGBA(image.Rect(0, 0, 1, 1))
	input.SetRGBA(0, 0, color.RGBA{R: 0x11, G: 0x22, B: 0x33, A: 0xFF})

	out := NewPaletteRecolorer().Recolor(input, assets, "skin", "unknown")
	got := color.RGBAModel.Convert(out.At(0, 0)).(color.RGBA)
	want := color.RGBA{R: 0x11, G: 0x22, B: 0x33, A: 0xFF}
	if got != want {
		t.Fatalf("expected unchanged pixel for unknown variant: got=%+v want=%+v", got, want)
	}
}

func TestPaletteRecolorerFuzzyVariantMatch(t *testing.T) {
	assets := t.TempDir()
	writePaletteFile(t, filepath.Join(assets, "palette_definitions", "skin_main.json"), `{"light":["#102030"],"auburn_dark":["#8090A0"]}`)

	input := image.NewRGBA(image.Rect(0, 0, 1, 1))
	input.SetRGBA(0, 0, color.RGBA{R: 0x10, G: 0x20, B: 0x30, A: 0xFF})

	out := NewPaletteRecolorer().Recolor(input, assets, "skin", "auburn")
	got := color.RGBAModel.Convert(out.At(0, 0)).(color.RGBA)
	want := color.RGBA{R: 0x80, G: 0x90, B: 0xA0, A: 0xFF}
	if got != want {
		t.Fatalf("unexpected fuzzy recolor result: got=%+v want=%+v", got, want)
	}
}

func TestPaletteRecolorerToleranceAndTransparencyParity(t *testing.T) {
	assets := t.TempDir()
	writePaletteFile(t, filepath.Join(assets, "palette_definitions", "skin_main.json"), `{"light":["#101010"],"tan":["#808080"]}`)

	input := image.NewRGBA(image.Rect(0, 0, 3, 1))
	input.SetRGBA(0, 0, color.RGBA{R: 0x12, G: 0x10, B: 0x0F, A: 0xFF}) // within tolerance
	input.SetRGBA(1, 0, color.RGBA{R: 0x13, G: 0x10, B: 0x10, A: 0xFF}) // outside tolerance
	input.SetRGBA(2, 0, color.RGBA{R: 0x10, G: 0x10, B: 0x10, A: 0x00}) // transparent

	out := NewPaletteRecolorer().Recolor(input, assets, "skin", "tan")

	if got := color.RGBAModel.Convert(out.At(0, 0)).(color.RGBA); got != (color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}) {
		t.Fatalf("expected tolerance match recolor at x=0, got %+v", got)
	}
	if got := color.RGBAModel.Convert(out.At(1, 0)).(color.RGBA); got != (color.RGBA{R: 0x13, G: 0x10, B: 0x10, A: 0xFF}) {
		t.Fatalf("expected no recolor outside tolerance at x=1, got %+v", got)
	}
	if got := color.RGBAModel.Convert(out.At(2, 0)).(color.RGBA); got.A != 0 {
		t.Fatalf("expected transparent pixel unchanged at x=2, got %+v", got)
	}
}

func writePaletteFile(t *testing.T, path, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}
