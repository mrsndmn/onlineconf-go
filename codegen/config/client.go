package config

import (
	"path/filepath"

	"github.com/mrsndmn/onlineconf-go/expr"
	"goa.design/goa/v3/codegen"
)

const (
	configModuleStructName         = "Module"
	configModuleReloaderStructName = "ModuleReloader"
)

func ConfigFile(config *expr.ConfigExpr) *codegen.File {

	snakeConfigName := codegen.SnakeCase(config.Name)
	path := filepath.Join(codegen.Gendir, snakeConfigName, "config.go")

	var (
		sections []*codegen.SectionTemplate
	)
	{
		header := codegen.Header(config.Name+" config", snakeConfigName,
			[]*codegen.ImportSpec{
				{Path: "context"},
				{Path: "io"},
				codegen.GoaImport(""),
			})
		def := &codegen.SectionTemplate{
			Name:   "config-struct",
			Source: configModuleT,
			Data:   config,
		}
		init := &codegen.SectionTemplate{
			Name:   "config-init",
			Source: configModuleInitT,
			Data:   config,
		}
		sections = []*codegen.SectionTemplate{header, def, init}
		// for _, m := range data.Methods {
		// 	sections = append(sections, &codegen.SectionTemplate{
		// 		Name:   "client-method",
		// 		Source: serviceClientMethodT,
		// 		Data:   m,
		// 	})
		// }
	}

	return &codegen.File{Path: path, SectionTemplates: sections}
}

// input: configData
const configModuleT = `// {{ .ModuleVarName }} is the {{ printf "%q" .Name }} service client.
type {{ .ModuleVarName }} struct {
	prefix string
{{- range .Params}}
	{{ .Name }} {{ .DType.Name() }}
{{- end }}
}
`

// input: configData
const configModuleInitT = `{{ printf "New%s initializes a %q service client given the endpoints." .ModuleVarName .Name | comment }}
func New{{ .ModuleVarName }}() *{{ .ModuleVarName }} {
	return &{{ .ModuleVarName }}{
		prefix: {{ .Prefix }},
{{- range .Params }}
	{{ .ParamName }}: { .DefaultValue }
{{- end }}
	}
}
`
