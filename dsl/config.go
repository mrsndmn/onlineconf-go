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
	expr.Root.Config = expr.NewConfigExpr(OCprefix, fn)
	return expr.Root.Config
}
