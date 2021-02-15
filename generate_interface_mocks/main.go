package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	currentFolder := filepath.Base(cwd)
	files := make([]string, 0)
	if err := filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() != currentFolder {
			return filepath.SkipDir
		}
		if filepath.Ext(path) == ".go" {
			files = append(files, info.Name())
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	for _, filename := range files {
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			log.Fatal("Error parsing file: " + filename)
		}

		base := strings.TrimSuffix(filename, filepath.Ext(filename))
		ast.Inspect(node, func(n ast.Node) bool {
			switch t := n.(type) {
			case *ast.TypeSpec:
				if t.Name.IsExported() {
					switch t.Type.(type) {
					case *ast.InterfaceType:
						out, err := exec.Command("mockery", "--name", t.Name.Name, "--filename", fmt.Sprintf("%s_%s.go", base, t.Name.Name)).CombinedOutput()
						if err != nil {
							log.Fatal(err)
						}
						log.Println(string(out))
					}
				}
			}
			return true
		})
	}
}
