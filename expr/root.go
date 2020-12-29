package expr

import (
	"github.com/mrsndmn/onlineconf-go/eval"
)

// Root is the root object built by the DSL.
var Root = &RootExpr{Configs: []*ConfigExpr{}}

type (
	// RootExpr is the struct built by the DSL on process start.
	RootExpr struct {
		// API contains the API expression built by the DSL.
		Configs []*ConfigExpr
	}
)

// WalkSets returns the expressions in order of evaluation.
func (r *RootExpr) WalkSets(walk eval.SetWalker) {
	if r.Configs == nil {
		return
	}

	for _, cfg := range r.Configs {
		walk(eval.ExpressionSet{cfg})
	}
}

// DependsOn returns nil, the core DSL has no dependency.
func (r *RootExpr) DependsOn() []eval.Root { return nil }

// Packages returns the Go import path to this and the dsl packages.
func (r *RootExpr) Packages() []string {
	return []string{
		"github.com/mrsndmn/onlineconf-go/expr",
		"github.com/mrsndmn/onlineconf-go/dsl",
	}
}

// EvalName is the name of the DSL.
func (r *RootExpr) EvalName() string {
	return "design"
}

// Validate makes sure the root expression is valid for code generation.
func (r *RootExpr) Validate() error {
	var verr eval.ValidationErrors
	if len(r.Configs) == 0 {
		verr.Add(r, "Missing Configs declaration")
	}

	return &verr
}

// Finalize finalizes the server expressions.
func (r *RootExpr) Finalize() {

}
