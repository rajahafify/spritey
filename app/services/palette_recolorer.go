package services

import (
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const recolorTolerance = 2

var basePaletteVariants = []string{"light", "base", "beige", "brown"}

type PaletteRecolorer struct{}

func NewPaletteRecolorer() PaletteRecolorer {
	return PaletteRecolorer{}
}

func (recolorer PaletteRecolorer) Recolor(img image.Image, assetsPath string, material string, variant string) image.Image {
	if img == nil || material == "" || variant == "" {
		return img
	}

	sourcePalette, ok := recolorer.getBasePalette(assetsPath, material)
	if !ok {
		return img
	}
	targetPalette, ok := recolorer.loadPalette(assetsPath, material, variant)
	if !ok {
		return img
	}

	bounds := img.Bounds()
	result := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(result, result.Bounds(), img, bounds.Min, draw.Src)

	pairCount := len(sourcePalette)
	if len(targetPalette) < pairCount {
		pairCount = len(targetPalette)
	}
	if pairCount == 0 {
		return img
	}

	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			original := color.RGBAModel.Convert(img.At(bounds.Min.X+x, bounds.Min.Y+y)).(color.RGBA)
			if original.A == 0 {
				continue
			}

			offset := result.PixOffset(x, y)
			for i := 0; i < pairCount; i++ {
				source := sourcePalette[i]
				target := targetPalette[i]
				if withinTolerance(original.R, source.R) &&
					withinTolerance(original.G, source.G) &&
					withinTolerance(original.B, source.B) {
					result.Pix[offset+0] = target.R
					result.Pix[offset+1] = target.G
					result.Pix[offset+2] = target.B
				}
			}
		}
	}

	return result
}

func (recolorer PaletteRecolorer) getBasePalette(assetsPath string, material string) ([]color.RGBA, bool) {
	for _, variant := range basePaletteVariants {
		palette, ok := recolorer.loadPalette(assetsPath, material, variant)
		if ok {
			return palette, true
		}
	}
	return nil, false
}

func (recolorer PaletteRecolorer) loadPalette(assetsPath string, material string, variant string) ([]color.RGBA, bool) {
	paletteFiles := findPaletteFiles(assetsPath, material)
	for _, path := range paletteFiles {
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var raw map[string]interface{}
		if err := json.Unmarshal(content, &raw); err != nil {
			continue
		}

		if candidate, ok := parsePaletteVariant(raw[variant]); ok {
			return candidate, true
		}

		keys := make([]string, 0, len(raw))
		for key := range raw {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			if strings.Contains(strings.ToLower(key), strings.ToLower(variant)) ||
				strings.Contains(strings.ToLower(variant), strings.ToLower(key)) {
				if candidate, ok := parsePaletteVariant(raw[key]); ok {
					return candidate, true
				}
			}
		}
	}
	return nil, false
}

func findPaletteFiles(assetsPath string, material string) []string {
	root := filepath.Join(assetsPath, "palette_definitions")
	matched := []string{}
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d == nil || d.IsDir() {
			return nil
		}
		name := d.Name()
		if strings.HasPrefix(name, material+"_") && strings.HasSuffix(strings.ToLower(name), ".json") {
			matched = append(matched, path)
		}
		return nil
	})
	sort.Strings(matched)
	return matched
}

func parsePaletteVariant(raw interface{}) ([]color.RGBA, bool) {
	values, ok := raw.([]interface{})
	if !ok || len(values) == 0 {
		return nil, false
	}

	palette := make([]color.RGBA, 0, len(values))
	for _, value := range values {
		hex, ok := value.(string)
		if !ok {
			return nil, false
		}
		rgb, ok := hexToRGBA(hex)
		if !ok {
			return nil, false
		}
		palette = append(palette, rgb)
	}
	return palette, true
}

func hexToRGBA(hex string) (color.RGBA, bool) {
	trimmed := strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(trimmed) != 6 {
		return color.RGBA{}, false
	}
	var components [3]uint8
	for i := 0; i < 3; i++ {
		component := trimmed[i*2 : i*2+2]
		value, ok := parseHexByte(component)
		if !ok {
			return color.RGBA{}, false
		}
		components[i] = value
	}
	return color.RGBA{R: components[0], G: components[1], B: components[2], A: 0xFF}, true
}

func parseHexByte(value string) (uint8, bool) {
	const (
		zero = byte('0')
		nine = byte('9')
		a    = byte('a')
		f    = byte('f')
		A    = byte('A')
		F    = byte('F')
	)
	if len(value) != 2 {
		return 0, false
	}

	parseNibble := func(ch byte) (uint8, bool) {
		switch {
		case ch >= zero && ch <= nine:
			return ch - zero, true
		case ch >= a && ch <= f:
			return ch - a + 10, true
		case ch >= A && ch <= F:
			return ch - A + 10, true
		default:
			return 0, false
		}
	}

	high, ok := parseNibble(value[0])
	if !ok {
		return 0, false
	}
	low, ok := parseNibble(value[1])
	if !ok {
		return 0, false
	}
	return (high << 4) | low, true
}

func withinTolerance(value uint8, source uint8) bool {
	delta := int(value) - int(source)
	if delta < 0 {
		delta = -delta
	}
	return delta <= recolorTolerance
}
