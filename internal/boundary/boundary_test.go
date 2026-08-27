package boundary

import (
	"encoding/json"
	"errors"
	"io"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var forbiddenOptkitImports = []string{
	"github.com/go-go-golems/coinvault",
	"github.com/go-go-golems/judgekit",
	"github.com/go-go-golems/ragkit",
	"github.com/go-go-golems/ragopt",
	"github.com/the-tree-center/rag-ttc",
}

type listedPackage struct {
	ImportPath   string
	Imports      []string
	TestImports  []string
	XTestImports []string
}

// TestOptkitRemainsDomainNeutral prevents product, RAG-domain, and
// measurement-domain packages from entering any Optkit package or test. Product
// integration belongs in the product repository and may import both sides.
func TestOptkitRemainsDomainNeutral(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve boundary test path")
	}
	moduleRoot := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", ".."))
	command := exec.Command("go", "list", "-json", "./...")
	command.Dir = moduleRoot
	output, err := command.Output()
	if err != nil {
		t.Fatalf("go list Optkit packages: %v", err)
	}
	decoder := json.NewDecoder(strings.NewReader(string(output)))
	checked := 0
	for {
		var pkg listedPackage
		if err := decoder.Decode(&pkg); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			t.Fatalf("decode go list package: %v", err)
		}
		checked++
		imports := append([]string{}, pkg.Imports...)
		imports = append(imports, pkg.TestImports...)
		imports = append(imports, pkg.XTestImports...)
		for _, imported := range imports {
			for _, forbidden := range forbiddenOptkitImports {
				if strings.HasPrefix(imported, forbidden) {
					t.Errorf("Optkit package %s imports forbidden domain package %s", pkg.ImportPath, imported)
				}
			}
		}
	}
	if checked == 0 {
		t.Fatal("go list returned no Optkit packages")
	}
}
