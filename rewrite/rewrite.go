//go:build ignore

// rewrite/rewrite.go
package rewrite

import (
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}

// Config holds the configuration for the import rewrite
type Config struct {
	OldPath string
	NewPath string
}

// Run executes the import rewrite with the given configuration
func Run() error {
	cfg := Config{
		OldPath: "github.com/JAD-RAD/aurabase",
		NewPath: "github.com/JAD-RAD/aurabase",
	}

	fmt.Printf("Rewriting imports from %s to %s\n", cfg.OldPath, cfg.NewPath)
	return RewriteImports(cfg)
}

// RewriteImports rewrites all import paths in Go files under the current directory
func RewriteImports(cfg Config) error {
	count := 0
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			if err := rewriteFile(path, cfg); err != nil {
				fmt.Printf("Error processing %s: %v\n", path, err)
			} else {
				count++
				fmt.Printf("Processed: %s\n", path)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking directory: %w", err)
	}

	fmt.Printf("Successfully processed %d files\n", count)
	return nil
}

func rewriteFile(filename string, cfg Config) error {
	// Read the file
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Parse the source file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error parsing file: %w", err)
	}

	// Update imports
	modified := false
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, "\"")
		if strings.HasPrefix(path, cfg.OldPath) {
			newPath := strings.Replace(path, cfg.OldPath, cfg.NewPath, 1)
			imp.Path.Value = fmt.Sprintf("%q", newPath)
			modified = true
		}
	}

	if !modified {
		return nil
	}

	// Format the file
	var buf []byte
	buf, err = format.Node(file, fset, file)
	if err != nil {
		return fmt.Errorf("error formatting file: %w", err)
	}

	// Write the changes back
	return ioutil.WriteFile(filename, buf, 0644)
}
