package dsl

import (
	"github.com/mrsndmn/onlineconf-go/eval"
	"github.com/mrsndmn/onlineconf-go/expr"
)

func Config(prefix string, fn func()) *expr.ConfigExpr {

	OCprefix := expr.OnlineConfPath(prefix)
	if !OCprefix.IsValid() {
		eval.ReportError("Config path is not valid OnlineConf prefix")
		return nil
	}

	if _, ok := eval.Current().(eval.TopExpr); !ok {
		eval.IncompatibleDSL()
		return nil
	}

	confExpr := expr.NewConfigExpr(OCprefix, fn)
	expr.Root.Configs = append(expr.Root.Configs, confExpr)

	eval.Execute(fn, confExpr)

	return confExpr
}

func SubConfig(prefix string, fn func()) *expr.ConfigExpr {

	OCprefix := expr.OnlineConfPath(prefix)
	if !OCprefix.IsValid() {
		eval.ReportError("SubConfig path is not valid OnlineConf prefix")
		return nil
	}

	topConfig, ok := eval.Current().(*expr.ConfigExpr)
	if !ok {
		eval.IncompatibleDSL()
		return nil
	}

	subCconfExpr := expr.NewConfigExpr(topConfig.Prefix+OCprefix, fn)
	expr.Root.Configs = append(expr.Root.Configs, subCconfExpr)

	eval.Execute(fn, subCconfExpr)

	return subCconfExpr
}
