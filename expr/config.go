package expr

import (
	"goa.design/goa/v3/eval"
)

type (
	// ConfigExpr contains the global properties for a API expression.
	ConfigExpr struct {
		// DSLFunc contains the DSL used to initialize the expression.
		eval.DSLFunc
		// Name of API
		Name   string
		Prefix OnlineConfPath

		ParamsByOCPath map[OnlineConfPath]*ParamExpr // params by onlineconf path
	}
)

// NewConfigExpr initializes an API expression.
func NewConfigExpr(name OnlineConfPath, dsl func()) *ConfigExpr {
	// todo
	return &ConfigExpr{
		ParamsByOCPath: map[OnlineConfPath]*ParamExpr{},
		DSLFunc:        dsl,
	}
}

// EvalName is the qualified name of the expression.
func (c *ConfigExpr) EvalName() string { return "Config " + c.Name }

// Hash returns a unique hash value for c.
func (c *ConfigExpr) Hash() string { return "_config_+" + c.Name }

// Finalize makes sure that the API name is initialized and there is at least
// one server definition (if none exists, it creates a default server). If API
// name is empty, it sets the name of the first service definition as API name.
func (c *ConfigExpr) Finalize() {
	// todo check all params has uniq names and uniq paths
}
