package generator

import (
	"fmt"

	cgconfig "github.com/mrsndmn/onlineconf-go/codegen/config"
	"github.com/mrsndmn/onlineconf-go/expr"
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
)

func Config(genpkg string, roots []eval.Root) ([]*codegen.File, error) {
	var files []*codegen.File
	for _, root := range roots {
		switch r := root.(type) {
		case *expr.RootExpr:
			for _, conf := range r.Configs {
				configFile := cgconfig.ConfigFile(conf)
				files = append(files, configFile)
			}
		}
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("design must define at least one service")
	}
	return files, nil
}
