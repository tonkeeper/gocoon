package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func abs(path string) string {
	p, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("abs path: %s", err)
	}
	return p
}

func main() {
	schemaFile := flag.String("schema", "cocoon_api.tlo", "path to compiled .tlo schema file")
	typesOutputDir := flag.String("typesOutputDir", "types/", "types output directory (will be wiped)")
	functionsOutputDir := flag.String("functionsOutputDir", "./", "functions output directory (written, not wiped)")
	typesPackageName := flag.String("typesPackage", "tlcocoonTypes", "generated types package name")
	functionsPackageName := flag.String("functionsPackage", "tlcocoon", "generated functions package name")
	typesImportPath := flag.String("typesImportPath", "github.com/tonkeeper/gocoon/pkg/tlcocoon/types", "import path for generated types package")
	flag.Parse()

	data, err := os.ReadFile(*schemaFile)
	if err != nil {
		log.Fatalf("read schema: %s", err)
	}

	s, err := ParseTLO(data)
	if err != nil {
		log.Fatalf("parse TLO: %s", err)
	}
	log.Printf("parsed TLO v%d: %d types, %d constructors, %d functions",
		s.Version, len(s.Types), len(s.Constructors), len(s.Functions))

	typesCodeByFile, functionsCodeByFile, err := generate(
		s,
		*typesPackageName,
		*functionsPackageName,
		*typesImportPath,
	)
	if err != nil {
		log.Fatalf("generate: %s", err)
	}

	typesDir := abs(*typesOutputDir)
	if err := os.RemoveAll(typesDir); err != nil {
		log.Fatalf("remove output dir: %s", err)
	}
	if err := os.MkdirAll(typesDir, 0o755); err != nil {
		log.Fatalf("mkdir: %s", err)
	}
	writeFiles(typesDir, typesCodeByFile)

	funcsDir := abs(*functionsOutputDir)
	if err := os.MkdirAll(funcsDir, 0o755); err != nil {
		log.Fatalf("mkdir functions dir: %s", err)
	}
	writeFiles(funcsDir, functionsCodeByFile)
}

func writeFiles(outDir string, codeByFile map[string]string) {
	names := make([]string, 0, len(codeByFile))
	for name := range codeByFile {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		outFile := filepath.Join(outDir, name)
		if err := os.WriteFile(outFile, []byte(codeByFile[name]), 0o644); err != nil {
			log.Fatalf("write %s: %s", outFile, err)
		}
		log.Printf("written %s", outFile)
	}
}
