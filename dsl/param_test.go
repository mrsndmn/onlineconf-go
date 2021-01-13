package dsl_test

import (
	"github.com/mrsndmn/onlineconf-go/dsl"
	"github.com/mrsndmn/onlineconf-go/expr"
	"goa.design/goa/v3/eval"
)

func (suite *CoonfigTestSuite) TestParamDSL() {
	dsl.Config("", func() {
		dsl.Param("/string", expr.String, func() {
			dsl.Required()
			dsl.Default("default string param value")
		})
		dsl.Param("/int", expr.Int, func() {
			dsl.Required()
			dsl.Default(100500)
		})
		dsl.Param("/float", expr.Float32, func() {
			dsl.Required()
			dsl.Default(3.1415)
		})
	})

	suite.Assert().Nil(eval.Context.Errors)
}

func (suite *CoonfigTestSuite) TestParamInvalidDefaultDSL() {
	testCases := map[string]eval.DSLFunc{
		"string invalid default": func() {
			dsl.Config("", func() {
				dsl.Param("/string", expr.String, func() {
					dsl.Default("default string param value")
				})
			})
		},
		"int invalid default": func() {
			dsl.Config("", func() {
				dsl.Param("/test_int", expr.Int, func() {
					dsl.Default(100500)
				})
			})
		},
		"uint invalid default": func() {
			dsl.Config("", func() {
				dsl.Param("/test_int", expr.UInt, func() {
					dsl.Default(-1)
				})
			})
		},
		"float invalid default": func() {
			dsl.Config("", func() {
				dsl.Param("/float", expr.Float32, func() {
					dsl.Default("")
				})
			})
		},
	}

	for k, dslFunc := range testCases {
		eval.Reset()
		dslFunc()
		suite.Assert().NotNil(eval.Context.Error(), k)
	}
}
