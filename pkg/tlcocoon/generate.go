//go:generate go run gen/main.go gen/codegen.go gen/tlo.go -schema cocoon_api.tlo -typesOutputDir types/ -functionsOutputDir . -typesPackage tlcocoonTypes -functionsPackage tlcocoon -typesImportPath github.com/tonkeeper/gocoon/pkg/tlcocoon/types
package tlcocoon
