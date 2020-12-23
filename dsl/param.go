package dsl

import (
	"github.com/mrsndmn/onlineconf-go/eval"
	"github.com/mrsndmn/onlineconf-go/expr"
)

func Param(paramPath string, dtype expr.DataType, dsl func()) {
	OCParamPath := expr.OnlineConfPath(paramPath)
	if !OCParamPath.IsValid() {
		eval.ReportError("Param path is not valid OnlineConf path")
		return
	}

	cfg, ok := eval.Current().(*expr.ConfigExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}

	if _, ok := cfg.ParamsByOCPath[OCParamPath]; ok {
		eval.ReportError("Parameter `%s` was already declared for config `%s`", OCParamPath, cfg.Prefix)
		return
	}

	paramExpr := expr.NewParamExpr(OCParamPath, dtype, dsl)

	cfg.ParamsByOCPath[OCParamPath] = paramExpr

	if !eval.Execute(dsl, paramExpr) {
		return
	}

	return
}

func Default(defaultValue interface{}) {
	param, ok := eval.Current().(*expr.ParamExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}

	// todo validate default value
	param.Default = defaultValue

	return
}

func Required() {
	param, ok := eval.Current().(*expr.ParamExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}

	param.Required = true

	return
}
