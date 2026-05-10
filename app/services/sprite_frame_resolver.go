package services

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

type SpriteFrameResolver struct{}

func NewSpriteFrameResolver() SpriteFrameResolver {
	return SpriteFrameResolver{}
}

var weaponAnimationPathMap = map[string]string{
	"slash":  "attack_slash",
	"thrust": "attack_thrust",
}

func (resolver SpriteFrameResolver) ResolveFrame(assetsPath string, pathPrefix string, animationID string) (string, bool, error) {
	prefix := cleanSpritePrefix(pathPrefix)

	structureA := path.Join(prefix, animationID+".png")
	if ok, err := fileExists(filepath.Join(assetsPath, "spritesheets", filepath.FromSlash(structureA))); err != nil {
		return "", false, err
	} else if ok {
		return structureA, true, nil
	}

	structureB := path.Join(prefix, "fg", animationID+".png")
	if ok, err := fileExists(filepath.Join(assetsPath, "spritesheets", filepath.FromSlash(structureB))); err != nil {
		return "", false, err
	} else if ok {
		return structureB, true, nil
	}

	structureCDir := path.Join(prefix, animationID)
	if frame, found, err := firstPNGInDirectory(filepath.Join(assetsPath, "spritesheets", filepath.FromSlash(structureCDir)), structureCDir, false); err != nil {
		return "", false, err
	} else if found {
		return frame, true, nil
	}

	mappedAnimationID, mapped := weaponAnimationPathMap[animationID]
	if mapped {
		structureDDir := path.Join(prefix, mappedAnimationID)
		if frame, found, err := firstPNGInDirectory(filepath.Join(assetsPath, "spritesheets", filepath.FromSlash(structureDDir)), structureDDir, true); err != nil {
			return "", false, err
		} else if found {
			return frame, true, nil
		}
	}

	return "", false, nil
}

func cleanSpritePrefix(raw string) string {
	normalized := strings.ReplaceAll(strings.TrimSpace(raw), "\\", "/")
	trimmed := strings.Trim(normalized, "/")
	if trimmed == "" {
		return ""
	}
	return path.Clean(trimmed)
}

func fileExists(filePath string) (bool, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !info.IsDir(), nil
}

func firstPNGInDirectory(dirPath string, dirPathRelative string, excludeBehind bool) (string, bool, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.ToLower(filepath.Ext(name)) != ".png" {
			continue
		}
		if excludeBehind && strings.Contains(name, "behind") {
			continue
		}
		return path.Join(dirPathRelative, name), true, nil
	}

	return "", false, nil
}
