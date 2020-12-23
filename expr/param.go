package expr

import (
	"github.com/mrsndmn/onlineconf-go/eval"
)

type (
	// ParamExpr contains the global properties for a API expression.
	ParamExpr struct {
		// DSLFunc contains the DSL used to initialize the expression.
		eval.DSLFunc
		// Name of API
		Name     string
		Path     OnlineConfPath
		DType    DataType
		Default  interface{}
		Required bool
	}
)

// NewParamExpr initializes an API expression.
func NewParamExpr(path OnlineConfPath, dtype DataType, dsl func()) *ParamExpr {
	curParam := &ParamExpr{
		Path:    path,
		DSLFunc: dsl,
		DType:   dtype,
	}

	return curParam
}

// EvalName is the qualified name of the expression.
func (c *ParamExpr) EvalName() string { return "Param " + c.Name }

// Hash returns a unique hash value for c.
func (c *ParamExpr) Hash() string { return "_param_+" + c.Name }

// Finalize makes sure that the API name is initialized and there is at least
// one server definition (if none exists, it creates a default server). If API
// name is empty, it sets the name of the first service definition as API name.
func (c *ParamExpr) Finalize() {
}
