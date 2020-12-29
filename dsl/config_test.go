package dsl_test

import (
	"testing"

	"github.com/mrsndmn/onlineconf-go/dsl"
	"github.com/mrsndmn/onlineconf-go/eval"
	"github.com/mrsndmn/onlineconf-go/expr"

	"github.com/stretchr/testify/suite"
)

type CoonfigTestSuite struct {
	suite.Suite
}

func (suite *CoonfigTestSuite) SetupTest() {
	eval.Reset()
}

func (suite *CoonfigTestSuite) TestValidConfigPrefix() {
	testCases := map[string]eval.DSLFunc{
		"empty":  func() { dsl.Config("", func() {}) },
		"simple": func() { dsl.Config("/onlineconf-go", func() {}) },
		"nasted": func() { dsl.Config("/onlineconf-go/prefix", func() {}) },
	}

	for k, dslFunc := range testCases {
		eval.Reset()
		dslFunc()
		suite.Assert().Nil(eval.Context.Errors, k)
	}
}

func (suite *CoonfigTestSuite) TestInvalidConfigPrefix() {
	testCases := map[string]eval.DSLFunc{
		"invalid one slash path":   func() { dsl.Config("/", func() {}) },
		"prefix with / at the end": func() { dsl.Config("/onlineconf-go/", func() {}) },
		"double slash in prefix":   func() { dsl.Config("/onlineconf-go//prefix/", func() {}) },
	}

	for k, dslFunc := range testCases {
		eval.Reset()
		dslFunc()
		suite.Assert().NotNil(eval.Context.Errors, k)
	}
}

func (suite *CoonfigTestSuite) TestConfigDSL() {
	dsl.Config("/onlineconf-go", func() {
		dsl.SubConfig("/db/meta", func() {
			dsl.Param("/host", expr.String, func() { dsl.Required() })
			dsl.Param("/pass", expr.String, func() { dsl.Required() })
			dsl.Param("/user", expr.String, func() { dsl.Required() })
		})

		dsl.Param("/host", expr.String, func() { dsl.Required() })
	})

	suite.Assert().Nil(eval.Context.Errors)
}

func (suite *CoonfigTestSuite) TestBadConfigDSL() {

	testCases := map[string]eval.DSLFunc{
		"prohibit nasted Configs": func() {
			dsl.Config("/onlineconf-go/db/meta", func() {
				dsl.Config("/", func() {})
			})
		},
		"prohibit Required() Inside Config": func() {
			dsl.Config("/onlineconf-go/db/meta", func() {
				dsl.Required()
			})
		},
		"prohibit Required() inside SubConfig": func() {
			dsl.Config("/onlineconf-go/db/meta", func() {
				dsl.SubConfig("/onlineconf-go/db/meta", func() {
					dsl.Required()
				})
			})
		},
	}

	for description, testCase := range testCases {
		eval.Reset()
		testCase()
		suite.Assert().NotNil(eval.Context.Errors, description)
	}
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(CoonfigTestSuite))
}
