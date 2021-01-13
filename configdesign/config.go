package configdesign

import (
	"github.com/mrsndmn/onlineconf-go/dsl"
	"github.com/mrsndmn/onlineconf-go/expr"
)

var _ = dsl.Config("/test/configdesign", func() {
	dsl.Param("/int_param", expr.Boolean, func() {
		dsl.Required()
	})
})
