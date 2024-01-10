package core

import (
	"os"
	"path/filepath"
	"testing"
)

var testData = filepath.Join("..", "..", "testdata")

func TestInitCfg(t *testing.T) {
	cfg, err := NewConfig(&CLIFlags{})
	if err != nil {
		t.Fatal(err)
	}

	// In v3.0, these should have defaults.
	if cfg.StylesPath == "" {
		t.Fatal("StylesPath is empty")
	} else if len(cfg.Paths) == 0 {
		t.Fatal("Paths are empty")
	}

	if !IsDir(cfg.StylesPath) {
		t.Fatalf("%s is not a directory", cfg.StylesPath)
	}
}

func TestGetIgnores(t *testing.T) {
	found, err := IgnoreFiles(filepath.Join(testData, "styles"))
	if err != nil {
		t.Fatal(err)
	} else if len(found) != 2 {
		t.Fatalf("Expected 2 ignore files, found %d", len(found))
	}
}

func TestFindAsset(t *testing.T) {
	cfg, err := NewConfig(&CLIFlags{})
	if err != nil {
		t.Fatal(err)
	}
	cfg.StylesPath = filepath.Join(testData, "styles")
	cfg.Paths = append(cfg.Paths, cfg.StylesPath)

	found := FindAsset(cfg, "line.tmpl")
	if found == "" {
		t.Fatal("Expected to find line.tmpl")
	}

	found = FindAsset(cfg, "notfound")
	if found != "" {
		t.Fatal("Expected to not find notfound")
	}
}

func TestFindAssetDefault(t *testing.T) {
	cfg, err := NewConfig(&CLIFlags{})
	if err != nil {
		t.Fatal(err)
	}

	expected, err := DefaultStylesPath()
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(expected, "tmp.tmpl")

	err = os.WriteFile(target, []byte{}, os.ModePerm)
	if err != nil {
		t.Fatal("Failed to create file", err)
	}

	found := FindAsset(cfg, "tmp.tmpl")
	if found == "" {
		t.Fatal("Expected to find 'tmp.tmpl', got empty")
	}

	found = FindAsset(cfg, "notfound")
	if found != "" {
		t.Fatal("Expected to not find 'notfound', got", found)
	}

	err = os.Remove(target)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFallbackToDefault(t *testing.T) {
	cfg, err := NewConfig(&CLIFlags{})
	if err != nil {
		t.Fatal(err)
	}
	local := filepath.Join(testData, "styles")

	cfg.StylesPath = local
	cfg.Paths = append(cfg.Paths, local)

	expected, err := DefaultStylesPath()
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(expected, "tmp.tmpl")

	err = os.WriteFile(target, []byte{}, os.ModePerm)
	if err != nil {
		t.Fatal("Failed to create file", err)
	}

	found := FindAsset(cfg, "tmp.tmpl")
	if found == "" {
		t.Fatal("Expected to find 'tmp.tmpl', got empty", cfg.Paths)
	}

	err = os.Remove(target)
	if err != nil {
		t.Fatal(err)
	}
}