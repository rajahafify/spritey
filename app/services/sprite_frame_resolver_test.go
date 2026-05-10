package services

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSpriteFrameResolverStructureA(t *testing.T) {
	assets := t.TempDir()
	writeResolverFile(t, filepath.Join(assets, "spritesheets", "weapon", "sword", "male", "walk.png"))

	resolved, ok, err := NewSpriteFrameResolver().ResolveFrame(assets, "weapon/sword/male", "walk")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !ok {
		t.Fatal("expected frame resolution success")
	}
	if resolved != "weapon/sword/male/walk.png" {
		t.Fatalf("unexpected resolved path: %q", resolved)
	}
}

func TestSpriteFrameResolverStructureB(t *testing.T) {
	assets := t.TempDir()
	writeResolverFile(t, filepath.Join(assets, "spritesheets", "hair", "short", "fg", "walk.png"))

	resolved, ok, err := NewSpriteFrameResolver().ResolveFrame(assets, "hair/short", "walk")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !ok {
		t.Fatal("expected frame resolution success")
	}
	if resolved != "hair/short/fg/walk.png" {
		t.Fatalf("unexpected resolved path: %q", resolved)
	}
}

func TestSpriteFrameResolverStructureCFirstHit(t *testing.T) {
	assets := t.TempDir()
	writeResolverFile(t, filepath.Join(assets, "spritesheets", "effect", "sparkle", "walk", "b.png"))
	writeResolverFile(t, filepath.Join(assets, "spritesheets", "effect", "sparkle", "walk", "a.png"))

	resolved, ok, err := NewSpriteFrameResolver().ResolveFrame(assets, "effect/sparkle", "walk")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !ok {
		t.Fatal("expected frame resolution success")
	}
	if resolved != "effect/sparkle/walk/a.png" {
		t.Fatalf("unexpected resolved path: %q", resolved)
	}
}

func TestSpriteFrameResolverStructureDSlashAndThrust(t *testing.T) {
	cases := []struct {
		name    string
		animID  string
		mapped  string
		wantRel string
	}{
		{
			name:    "slash",
			animID:  "slash",
			mapped:  "attack_slash",
			wantRel: "weapon/spear/male/attack_slash/front.png",
		},
		{
			name:    "thrust",
			animID:  "thrust",
			mapped:  "attack_thrust",
			wantRel: "weapon/spear/male/attack_thrust/front.png",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assets := t.TempDir()
			writeResolverFile(t, filepath.Join(assets, "spritesheets", "weapon", "spear", "male", tc.mapped, "back_behind.png"))
			writeResolverFile(t, filepath.Join(assets, "spritesheets", "weapon", "spear", "male", tc.mapped, "front.png"))

			resolved, ok, err := NewSpriteFrameResolver().ResolveFrame(assets, "weapon/spear/male", tc.animID)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if !ok {
				t.Fatal("expected frame resolution success")
			}
			if resolved != tc.wantRel {
				t.Fatalf("unexpected resolved path: %q", resolved)
			}
		})
	}
}

func TestSpriteFrameResolverStructureCPrecedenceOverD(t *testing.T) {
	assets := t.TempDir()
	writeResolverFile(t, filepath.Join(assets, "spritesheets", "weapon", "axe", "male", "slash", "c-hit.png"))
	writeResolverFile(t, filepath.Join(assets, "spritesheets", "weapon", "axe", "male", "attack_slash", "d-hit.png"))

	resolved, ok, err := NewSpriteFrameResolver().ResolveFrame(assets, "weapon/axe/male", "slash")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !ok {
		t.Fatal("expected frame resolution success")
	}
	if resolved != "weapon/axe/male/slash/c-hit.png" {
		t.Fatalf("unexpected resolved path: %q", resolved)
	}
}

func TestSpriteFrameResolverNotFound(t *testing.T) {
	resolved, ok, err := NewSpriteFrameResolver().ResolveFrame(t.TempDir(), "weapon/sword/male", "slash")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ok {
		t.Fatalf("expected not found, got %q", resolved)
	}
}

func writeResolverFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte{0x89, 'P', 'N', 'G'}, 0o644); err != nil {
		t.Fatal(err)
	}
}
