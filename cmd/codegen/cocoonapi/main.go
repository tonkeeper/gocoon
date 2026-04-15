package main

import (
	"flag"
	"fmt"
	"go/format"
	"hash/crc32"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tonkeeper/tongo/tl/parser"
)

// normalizeDecl joins a (possibly multi-line) TL declaration into a single
// normalized string and strips the trailing semicolon, for CRC32 computation.
func normalizeDecl(decl string) string {
	lines := strings.Split(decl, "\n")
	var clean []string
	for _, l := range lines {
		if idx := strings.Index(l, "//"); idx >= 0 {
			l = l[:idx]
		}
		l = strings.TrimSpace(l)
		if l != "" {
			clean = append(clean, l)
		}
	}
	joined := strings.Join(clean, " ")
	spaceRe := regexp.MustCompile(`\s+`)
	joined = spaceRe.ReplaceAllString(joined, " ")
	return strings.TrimSuffix(strings.TrimSpace(joined), ";")
}

// addMissingTags inserts a CRC32-computed hex tag into any constructor that
// has no explicit tag, handling both single-line and multi-line declarations.
func addMissingTags(src string) string {
	lines := strings.Split(src, "\n")
	result := make([]string, 0, len(lines))

	i := 0
	for i < len(lines) {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Pass through blank lines, comments, and section markers.
		if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "---") {
			result = append(result, line)
			i++
			continue
		}

		// Indented lines are continuation lines, handled as part of their declaration.
		if len(line) > 0 && (line[0] == ' ' || line[0] == '\t') {
			result = append(result, line)
			i++
			continue
		}

		// Collect the full declaration (ends at the line containing ";").
		declLines := []string{line}
		j := i + 1
		if !strings.Contains(trimmed, ";") {
			for j < len(lines) {
				next := lines[j]
				declLines = append(declLines, next)
				j++
				if strings.Contains(strings.TrimSpace(next), ";") {
					break
				}
			}
		}

		// Not a combinator declaration if there's no "=" anywhere in it.
		fullDecl := strings.Join(declLines, "\n")
		if !strings.Contains(fullDecl, "=") {
			for _, l := range declLines {
				result = append(result, l)
			}
			i = j
			continue
		}

		constructorName := strings.FieldsFunc(trimmed, func(r rune) bool {
			return r == ' ' || r == '\t'
		})[0]

		// Already has an explicit tag.
		if strings.Contains(constructorName, "#") {
			for _, l := range declLines {
				result = append(result, l)
			}
			i = j
			continue
		}

		// Compute CRC32 tag from the normalized declaration.
		tag := fmt.Sprintf("#%08x", crc32.Checksum([]byte(normalizeDecl(fullDecl)), crc32.IEEETable))

		spaceIdx := strings.IndexAny(trimmed, " \t")
		var newFirstLine string
		if spaceIdx < 0 {
			newFirstLine = constructorName + tag
		} else {
			newFirstLine = constructorName + tag + trimmed[spaceIdx:]
		}
		result = append(result, newFirstLine)
		for _, l := range declLines[1:] {
			result = append(result, l)
		}
		i = j
	}

	return strings.Join(result, "\n")
}

func main() {
	input := flag.String("input", "pkg/tlcocoonapi/cocoon_api.tl", "path to cocoon_api.tl")
	outputDir := flag.String("output", "pkg/tlcocoonapi", "directory to write generated files")
	flag.Parse()

	scheme, err := os.ReadFile(*input)
	if err != nil {
		panic(err)
	}

	tagged := addMissingTags(string(scheme))

	parsed, err := parser.Parse(tagged)
	if err != nil {
		panic(err)
	}

	// Extend default known types with double (float64).
	knownTypes := map[string]parser.DefaultType{
		"#":      {"uint32", false},
		"int":    {"uint32", false},
		"int256": {"tl.Int256", false},
		"long":   {"uint64", false},
		"bytes":  {"[]byte", true},
		"Bool":   {"bool", false},
		"string": {"string", false},
		"double": {"float64", false},
	}

	g := parser.NewGenerator(knownTypes, "*Client")

	types, err := g.LoadTypes(parsed.Declarations)
	if err != nil {
		panic(err)
	}

	functions, err := g.LoadFunctionsWithTransport(parsed.Functions, "request")
	if err != nil {
		panic(err)
	}

	src := `// Code generated - DO NOT EDIT.

package tlcocoonapi

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/tonkeeper/tongo/tl"
	"io"
)
` + types + functions

	formatted, err := format.Source([]byte(src))
	if err != nil {
		fmt.Fprintln(os.Stderr, "format error:")
		fmt.Fprintln(os.Stderr, err)
		os.WriteFile(filepath.Join(*outputDir, "generated.go"), []byte(src), 0644)
		os.Exit(1)
	}
	outPath := filepath.Join(*outputDir, "generated.go")
	if err := os.WriteFile(outPath, formatted, 0644); err != nil {
		panic(err)
	}
	fmt.Println(outPath)
}
