package services

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestRecipeValidatorValidatesRecipe(t *testing.T) {
	assets := filepath.Join("..", "..", "testdata", "fixtures", "basic-assets")
	recipe := filepath.Join("..", "..", "testdata", "fixtures", "recipes", "valid-basic.json")

	result, problem := NewRecipeValidator().Validate(recipe, assets)
	if problem != nil {
		t.Fatalf("expected recipe to validate, got %v", problem)
	}

	if result.BodyType != "male" {
		t.Fatalf("unexpected body type: %q", result.BodyType)
	}
	if len(result.Selections) != 2 {
		t.Fatalf("expected 2 selections, got %+v", result.Selections)
	}
}

func TestRecipeValidatorUsesDefaultBodyType(t *testing.T) {
	assets := filepath.Join("..", "..", "testdata", "fixtures", "basic-assets")
	recipe := filepath.Join("..", "..", "testdata", "fixtures", "recipes", "default-body-type.json")

	result, problem := NewRecipeValidator().Validate(recipe, assets)
	if problem != nil {
		t.Fatalf("expected recipe to validate, got %v", problem)
	}
	if result.BodyType != "male" {
		t.Fatalf("expected default body type male, got %q", result.BodyType)
	}
}

func TestRecipeValidatorUnknownLayer(t *testing.T) {
	assets := filepath.Join("..", "..", "testdata", "fixtures", "basic-assets")
	recipe := filepath.Join("..", "..", "testdata", "fixtures", "recipes", "unknown-layer.json")

	_, problem := NewRecipeValidator().Validate(recipe, assets)
	if problem == nil {
		t.Fatal("expected validation problem")
	}
	if problem.Code != "UNKNOWN_LAYER_ID" {
		t.Fatalf("unexpected problem: %+v", problem)
	}
}

func TestRecipeValidatorMissingRequiredSpriteFrame(t *testing.T) {
	assets, recipe := writeReadinessValidationFixture(t, readinessFixtureOptions{
		packDefaultsAnimations: []string{"idle", "walk"},
		includeBodyIdle:        true,
		includeBodyWalk:        false,
		includeWeaponIdle:      true,
		includeWeaponWalk:      true,
	})

	_, problem := NewRecipeValidator().Validate(recipe, assets)
	if problem == nil {
		t.Fatal("expected validation problem")
	}
	if problem.Code != "MISSING_SPRITE_FRAME" {
		t.Fatalf("unexpected code: %+v", problem)
	}
}

func TestRecipeValidatorFallbackBodyTypePathAddsWarning(t *testing.T) {
	assets, recipe := writeReadinessValidationFixture(t, readinessFixtureOptions{
		packDefaultsAnimations: []string{"idle"},
		recipeBodyType:         "female",
		missingFallback:        "male",
		includeBodyIdle:        true,
		includeBodyFemaleIdle:  true,
		includeWeaponIdle:      true,
		weaponHasFemalePath:    false,
	})

	result, problem := NewRecipeValidator().Validate(recipe, assets)
	if problem != nil {
		t.Fatalf("expected recipe to validate, got %+v", problem)
	}
	if len(result.Warnings) != 1 {
		t.Fatalf("expected one warning, got %+v", result.Warnings)
	}
}

type readinessFixtureOptions struct {
	packDefaultsAnimations []string
	recipeBodyType         string
	missingFallback        string
	includeBodyIdle        bool
	includeBodyFemaleIdle  bool
	includeBodyWalk        bool
	includeWeaponIdle      bool
	includeWeaponWalk      bool
	weaponHasFemalePath    bool
}

func writeReadinessValidationFixture(t *testing.T, opts readinessFixtureOptions) (string, string) {
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
	weaponFemale := `"female":"weapon/sword/female/"`
	if !opts.weaponHasFemalePath {
		weaponFemale = ""
	}

	writeFixtureText(t, filepath.Join(assets, "pack.json"), `{"schema_version":"1","id":"readiness-assets","name":"Readiness Assets","defaults":{"body_type":"male","animations":["idle","walk"],"canvas_width":8,"missing_body_type_fallback":"`+missingFallback+`"}}`)
	if len(opts.packDefaultsAnimations) > 0 {
		animationsJSON := `["` + opts.packDefaultsAnimations[0] + `"`
		for i := 1; i < len(opts.packDefaultsAnimations); i++ {
			animationsJSON += `,"` + opts.packDefaultsAnimations[i] + `"`
		}
		animationsJSON += `]`
		writeFixtureText(t, filepath.Join(assets, "pack.json"), `{"schema_version":"1","id":"readiness-assets","name":"Readiness Assets","defaults":{"body_type":"male","animations":`+animationsJSON+`,"canvas_width":8,"missing_body_type_fallback":"`+missingFallback+`"}}`)
	}
	writeFixtureText(t, filepath.Join(assets, "sheet_definitions", "body", "body_human.json"), `{"name":"Human Body","type_name":"body","layer_1":{"zPos":10,"male":"body/human/male/","female":"body/human/female/"},"animations":["idle","walk"]}`)
	writeFixtureText(t, filepath.Join(assets, "sheet_definitions", "weapon", "sword_training.json"), `{"name":"Training Sword","type_name":"weapon","layer_1":{"zPos":30,"male":"weapon/sword/male/"`+withLeadingComma(weaponFemale)+`},"animations":["idle","walk"]}`)
	writeFixtureText(t, filepath.Join(assets, "palette_definitions", ".gitkeep"), "")
	writeFixtureText(t, filepath.Join(assets, "spritesheets", ".gitkeep"), "")

	if opts.includeBodyIdle {
		writeFixturePNG(t, filepath.Join(assets, "spritesheets", "body", "human", "male", "idle.png"), color.RGBA{R: 30, G: 80, B: 190, A: 255})
	}
	if opts.includeBodyWalk {
		writeFixturePNG(t, filepath.Join(assets, "spritesheets", "body", "human", "male", "walk.png"), color.RGBA{R: 35, G: 90, B: 200, A: 255})
	}
	if opts.includeBodyFemaleIdle {
		writeFixturePNG(t, filepath.Join(assets, "spritesheets", "body", "human", "female", "idle.png"), color.RGBA{R: 38, G: 95, B: 205, A: 255})
	}
	if opts.includeWeaponIdle {
		writeFixturePNG(t, filepath.Join(assets, "spritesheets", "weapon", "sword", "male", "idle.png"), color.RGBA{R: 190, G: 50, B: 80, A: 255})
	}
	if opts.includeWeaponWalk {
		writeFixturePNG(t, filepath.Join(assets, "spritesheets", "weapon", "sword", "male", "walk.png"), color.RGBA{R: 200, G: 60, B: 90, A: 255})
	}

	recipe := filepath.Join(root, "recipe.json")
	writeFixtureText(t, recipe, `{"body_type":"`+recipeBodyType+`","selections":{"body":{"id":"body_human"},"weapon":{"id":"sword_training"}}}`)
	return assets, recipe
}

func withLeadingComma(value string) string {
	if value == "" {
		return ""
	}
	return "," + value
}

func writeFixtureText(t *testing.T, path, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeFixturePNG(t *testing.T, path string, fill color.RGBA) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			c := fill
			c.B = uint8((int(c.B) + x + y) % 255)
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
