package eval

import "fmt"

var Context *DSLContext

type (
	// DSLContext is the data structure that contains the DSL execution state.
	DSLContext struct {
		// Stack represents the current execution stack.
		Stack Stack
		// Errors contains the DSL execution errors for the current expression set.
		// Errors is an instance of MultiError.
		Errors error

		// roots is the list of DSL roots as registered by all loaded DSLs.
		roots []Root
		// dslPackages keeps track of the DSL package import paths so the initiator
		// may skip any callstack frame that belongs to them when computing error
		// locations.
		dslPackages []string
	}

	// Stack represents the expression evaluation stack. The stack is appended to
	// each time the initiator executes an expression source DSL.
	Stack []Expression
)

func init() {
	Reset()
}

// Reset resets the eval context, mostly useful for tests.
func Reset() {
	Context = &DSLContext{dslPackages: []string{"goa.design/goa/v3/eval"}}
}

// Current returns current evaluation context, i.e. object being currently built
// by DSL.
func (s Stack) Current() Expression {
	if len(s) == 0 {
		return nil
	}
	return s[len(s)-1]
}

// Error builds the error message from the current context errors.
func (c *DSLContext) Error() string {
	if c.Errors != nil {
		return c.Errors.Error()
	}
	return ""
}

// Register appends a root expression to the current Context root expressions.
// Each root expression may only be registered once.
func Register(r Root) error {
	for _, o := range Context.roots {
		if r.EvalName() == o.EvalName() {
			return fmt.Errorf("duplicate DSL %s", r.EvalName())
		}
	}
	Context.dslPackages = append(Context.dslPackages, r.Packages()...)
	Context.roots = append(Context.roots, r)

	return nil
}

// Record appends an error to the context Errors field.
func (c *DSLContext) Record(err *Error) {
	if c.Errors == nil {
		c.Errors = MultiError{err}
	} else {
		c.Errors = append(c.Errors.(MultiError), err)
	}
}
