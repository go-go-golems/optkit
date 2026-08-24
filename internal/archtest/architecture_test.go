package archtest

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
)

const modulePath = "github.com/go-go-golems/optkit"

func TestProductionImportBoundaries(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve architecture test location")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
	var violations []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			name := entry.Name()
			if name == ".git" || name == "vendor" || name == "tmp" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		family := strings.Split(filepath.ToSlash(relative), "/")[0]
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, spec := range file.Imports {
			importPath, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				return err
			}
			if !strings.HasPrefix(importPath, modulePath+"/") {
				continue
			}
			importedRelative := strings.TrimPrefix(importPath, modulePath+"/")
			importedFamily := strings.Split(importedRelative, "/")[0]
			if reason := forbiddenImport(family, importedFamily, importedRelative); reason != "" {
				violations = append(violations, fmt.Sprintf("%s imports %s: %s", filepath.ToSlash(relative), importPath, reason))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(violations)
	if len(violations) > 0 {
		t.Fatalf("package boundary violations:\n  %s", strings.Join(violations, "\n  "))
	}
}

func forbiddenImport(family, importedFamily, importedRelative string) string {
	switch family {
	case "record":
		return "record must import only the standard library"
	case "artifact":
		if importedFamily != "record" && importedFamily != "artifact" {
			return "artifact may import only record and its own package family"
		}
	}

	core := map[string]bool{
		"space": true, "episode": true, "measure": true, "experiment": true,
		"campaign": true, "scheduler": true, "budget": true,
	}
	if core[family] {
		switch importedFamily {
		case "local", "store", "projection", "examples", "cmd":
			return "core packages cannot depend on concrete stores, projections, examples, or commands"
		}
	}
	if family == "measure" {
		switch importedFamily {
		case "campaign", "scheduler":
			return "measurement semantics cannot depend on campaign control or scheduling"
		}
	}
	if family == "projection" && (importedFamily == "store" || importedFamily == "local") {
		return "projections must be rebuildable without a concrete store dependency"
	}
	if strings.HasPrefix(importedRelative, "store/sqlite") {
		switch family {
		case "store", "local", "cmd", "examples":
			return ""
		default:
			return "only composition roots may import the concrete SQLite store"
		}
	}
	return ""
}
